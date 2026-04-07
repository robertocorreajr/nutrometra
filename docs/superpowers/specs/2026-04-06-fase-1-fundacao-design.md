# Spec: Fase 1 — Fundação

**Data:** 2026-04-06
**Escopo:** Monorepo scaffold + ambiente Docker Compose + Go API core (auth, multi-tenant, RBAC, planos/entitlements)
**Fase do MVP:** Fase 1 de 5

---

## 1. Contexto e suposições

- Nenhum código existe ainda — repositório contém apenas documentação.
- Backend principal em Go conforme definido no CLAUDE.md.
- Provedor de autenticação: **Zitadel** self-hosted (Go-native, API-first, multi-tenancy built-in).
- Banco de dados principal: **PostgreSQL 16**.
- Cache e filas leves: **Redis 7**.
- Frontends (Next.js) ficam como scaffold vazio nesta fase — nenhuma lógica de UI implementada.
- Go worker fica como scaffold vazio nesta fase.
- Ambiente local 100% reproduzível via Docker Compose sem configuração manual.
- Zitadel usa database separado (`zitadel`) na mesma instância PostgreSQL do ambiente local. Em produção, instâncias separadas.
- Segredos carregados via variáveis de ambiente; `.env` no `.gitignore`.

---

## 2. Estrutura do monorepo

```
nutrometra/
├── apps/
│   ├── web-professional/        # Next.js — scaffold (package.json, tsconfig)
│   ├── web-patient/             # Next.js — scaffold
│   └── backoffice/              # Next.js — scaffold
├── services/
│   ├── api/                     # Go API — entregável principal desta fase
│   └── worker/                  # Go worker — scaffold
├── infra/
│   ├── docker-compose.yml       # ambiente local completo
│   ├── docker-compose.override.yml.example
│   ├── postgres/
│   │   └── init.sql             # cria databases: nutrometra, zitadel
│   └── zitadel/
│       └── config.yaml          # config declarativa: realm, client, admin user
├── docs/
├── agents/
├── Makefile                     # targets: dev, test, migrate, lint, build
├── .env.example
└── CLAUDE.md
```

---

## 3. Docker Compose — serviços

| Serviço    | Imagem                          | Porta | Observação                                    |
|------------|---------------------------------|-------|-----------------------------------------------|
| `postgres`  | `postgres:16-alpine`            | 5432  | dois databases: `nutrometra` e `zitadel`      |
| `redis`     | `redis:7-alpine`                | 6379  | cache e sessões                               |
| `zitadel`   | `ghcr.io/zitadel/zitadel:latest`| 8080  | modo `start-from-init`, config via YAML       |
| `api`       | build local Go                  | 8081  | hot-reload via `air`                         |

- `zitadel` depende de `postgres` com healthcheck.
- `api` depende de `postgres`, `redis` e `zitadel` com healthcheck.
- Volumes nomeados para persistência de dados entre restarts (`postgres_data`, `redis_data`).
- `.env.example` documenta todas as variáveis necessárias.

---

## 4. Estrutura do Go API (`services/api/`)

```
services/api/
├── cmd/api/
│   └── main.go                  # bootstrap: config → db → redis → server
├── internal/
│   ├── platform/
│   │   ├── config/              # carrega e valida env vars (struct tipada)
│   │   ├── db/                  # pool pgxpool, helper RunInTx
│   │   ├── redis/               # cliente go-redis
│   │   ├── logger/              # slog: JSON (prod), texto (dev)
│   │   ├── server/              # chi router, graceful shutdown
│   │   └── observability/       # /health, /ready, /metrics (Prometheus)
│   ├── identity/
│   │   ├── domain/              # User, Claims, Session
│   │   ├── usecase/             # ValidateToken, GetMe, Logout
│   │   └── oidc/                # adapter Zitadel: JWKS cache, token validation
│   ├── tenancy/
│   │   ├── domain/              # Tenant, TenantUser, TenantStatus
│   │   ├── usecase/             # ProvisionTenant, ResolveTenant, SuspendTenant
│   │   └── repository/          # postgres: tenants, tenant_users
│   ├── billing/
│   │   ├── domain/              # Plan, Subscription, Feature, Entitlement
│   │   ├── usecase/             # CheckEntitlement, ActivateTrial, GetSubscription
│   │   └── repository/          # postgres: subscription_plans, tenant_subscriptions, ...
│   └── rbac/
│       ├── domain/              # Role, Permission, Policy
│       ├── usecase/             # Enforce, AssignRole, RevokeRole
│       └── repository/          # postgres: roles, permissions, tenant_user_roles
├── migrations/                  # golang-migrate, arquivos .up.sql e .down.sql
├── go.mod
├── go.sum
├── .air.toml                    # hot-reload config
└── Makefile
```

---

## 5. Middleware stack

Todo request passa pela cadeia na ordem abaixo:

```
HTTP request
  → RequestID         (gera/propaga X-Request-ID)
  → Logger            (loga método, path, status, latência)
  → Recover           (captura panics, retorna 500)
  → AuthMiddleware    (valida JWT Zitadel via JWKS; 401 se ausente/inválido)
  → TenantMiddleware  (resolve tenant_id do claim; 403 se suspenso)
  → RBACMiddleware    (verifica permissão requerida; 403 se negado)
  → Handler
```

Endpoints públicos (`/health`, `/ready`, `/plans`, `/auth/callback`) bypassam `AuthMiddleware`.

**Contexto propagado em toda operação I/O:**
`request_id`, `tenant_id`, `user_id`, `actor_role`, `trace_id`

---

## 6. Migrations

15 migrations cobrindo o núcleo SaaS da Fase 1:

| # | Nome | Conteúdo |
|---|------|----------|
| 001 | `create_tenants` | tabela `tenants` |
| 002 | `create_users` | tabela `users` |
| 003 | `create_tenant_users` | tabela `tenant_users` |
| 004 | `create_roles_permissions` | `roles`, `permissions`, `role_permissions` |
| 005 | `create_tenant_user_roles` | tabela `tenant_user_roles` |
| 006 | `create_subscription_plans` | tabela `subscription_plans` |
| 007 | `create_plan_features` | tabela `plan_features` |
| 008 | `create_coupons` | tabela `coupons` |
| 009 | `create_tenant_subscriptions` | tabela `tenant_subscriptions` |
| 010 | `create_billing_invoices` | tabela `billing_invoices` |
| 011 | `create_billing_payments` | tabela `billing_payments` |
| 012 | `create_tenant_feature_overrides` | tabela `tenant_feature_overrides` |
| 013 | `create_audit_logs` | tabela `audit_logs` |
| 014 | `seed_roles_permissions` | roles: `owner`, `nutritionist`, `receptionist`, `patient`, `backoffice_admin`; permissões base |
| 015 | `seed_subscription_plans` | planos: `free`, `starter`, `pro`, `enterprise` com features iniciais |

**Convenções:**
- `tenant_id` obrigatório em todas as tabelas funcionais (exceto `users`, `roles`, `permissions`, `subscription_plans` — globais).
- `created_at`, `updated_at` em todas as tabelas mutáveis.
- Índice composto em `(tenant_id, status)` nas tabelas de alta leitura.
- Todas as PKs são UUIDs gerados pela aplicação.

---

## 7. Endpoints da Fase 1

### Autenticação
```
POST /auth/callback          # troca code Zitadel por sessão
POST /auth/refresh           # renova access token
POST /auth/logout            # invalida sessão
GET  /auth/me                # perfil + roles + tenant do usuário autenticado
```

### Tenants
```
POST   /tenants                          # provisiona tenant (signup)
GET    /tenants/:id                      # dados do tenant
PATCH  /tenants/:id                      # atualiza nome, timezone, locale
GET    /tenants/:id/members              # lista membros com roles
POST   /tenants/:id/members/invite       # convida membro por email
DELETE /tenants/:id/members/:user_id     # remove membro
```

### Planos e assinaturas
```
GET  /plans                              # lista planos ativos (público)
GET  /tenants/:id/subscription           # assinatura atual
POST /tenants/:id/subscription/trial     # ativa trial
GET  /tenants/:id/entitlements           # mapa de features habilitadas
```

### RBAC
```
GET    /roles                                        # lista roles disponíveis
POST   /tenants/:id/members/:uid/roles               # atribui role
DELETE /tenants/:id/members/:uid/roles/:role_id      # remove role
```

### Operacional
```
GET /health    # status: db, redis, zitadel
GET /ready     # readiness probe (k8s-ready)
GET /metrics   # Prometheus (requer auth interna)
```

**Contrato de erro padrão:**
```json
{
  "code": "permission_denied",
  "message": "You do not have permission to perform this action.",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 8. RBAC — roles e permissões iniciais

| Role | Escopo | Permissões representativas |
|------|--------|---------------------------|
| `owner` | tenant | `tenant:manage`, `members:invite`, `roles:assign`, `subscription:manage` |
| `nutritionist` | tenant | `patients:read`, `clinical:write`, `schedule:manage` |
| `receptionist` | tenant | `schedule:read`, `patients:read` |
| `patient` | tenant | `self:read`, `diet:read`, `appointments:read` |
| `backoffice_admin` | backoffice | `tenants:manage`, `billing:manage`, `support:manage` |

- Deny-by-default: sem role atribuída = zero acesso.
- Backoffice usa `actor_scope = backoffice`; nunca mistura com escopo de tenant.
- Permissões armazenadas como `code` string (ex: `clinical:write`), verificadas no middleware.

---

## 9. Segurança

- JWT RS256 emitido pelo Zitadel, validado via JWKS com cache local (TTL 5 min, refresh automático).
- `access_token` com vida de 15 min; `refresh_token` rotacionado a cada uso.
- Proteção CSRF via double-submit cookie em endpoints de mutação.
- MFA disponível via Zitadel, desabilitado por padrão.
- Nenhum segredo em código ou repositório. `.env` no `.gitignore`.
- Erros internos nunca expõem stack traces nem detalhes de infraestrutura ao cliente.
- Zitadel service account key para comunicação machine-to-machine (worker → Zitadel).

---

## 10. Auditoria

Toda mutação crítica grava em `audit_logs` na mesma transação da operação:

| Evento auditado | `entity_type` | `action` |
|-----------------|---------------|----------|
| Provisionar tenant | `tenant` | `created` |
| Suspender tenant | `tenant` | `suspended` |
| Atribuir role | `tenant_user_role` | `assigned` |
| Revogar role | `tenant_user_role` | `revoked` |
| Ativar trial | `tenant_subscription` | `trial_activated` |
| Cancelar assinatura | `tenant_subscription` | `cancelled` |
| Login | `user_session` | `login` |
| Logout | `user_session` | `logout` |

Campos obrigatórios: `tenant_id`, `actor_user_id`, `actor_scope`, `entity_type`, `entity_id`, `action`, `ip_address`, `created_at`.

---

## 11. Observabilidade

- Structured logging com `slog`: JSON em produção, texto legível em dev.
- `request_id` gerado no ingresso de cada request, propagado em contexto e retornado no header `X-Request-ID`.
- `/health` verifica conectividade ativa com PostgreSQL, Redis e Zitadel.
- `/ready` retorna 200 apenas quando todas as dependências estão saudáveis e migrations aplicadas.
- Prometheus endpoint em `/metrics` com métricas básicas: latência por endpoint, taxa de erro, conexões de DB.

---

## 12. Entitlement check — lógica de runtime

```
CheckEntitlement(ctx, tenant_id, feature_key) → (enabled bool, limit int, err)

  1. Carrega tenant_subscription ativo para o tenant
  2. Verifica tenant_feature_overrides (override tem precedência sobre plano)
  3. Verifica plan_features do plano ativo
  4. Retorna (enabled, limit)
```

Usado como guard em handlers e como middleware configurável por rota.

---

## 13. Riscos e pontos pendentes

| Risco | Mitigação |
|-------|-----------|
| Zitadel config YAML pode divergir entre ambientes | Config declarativa versionada no repo; script de init idempotente |
| JWKS cache stale durante rotation de chaves | TTL curto (5 min) + refresh em caso de 401 inesperado |
| Migrations em produção requerem downtime zero | Migrations só adicionivas no MVP; drops/renames via fases separadas |
| Frontends Next.js ficam sem auth nesta fase | Aceitável — Phase 1 entrega API testável; portais vêm na Fase 2 |
| Billing provider (Stripe etc.) não integrado ainda | Fase 4 — assinatura manual via backoffice até lá |

---

## 14. Notas para fases posteriores

- **CRN obrigatório no cadastro de nutricionista (Fase 2):** a tabela `professionals` deve registrar `registration_type = CRN`, `registration_number` (ex: `12345/P`) e `registration_state` (sigla do estado do conselho regional, ex: `SP`). O CRN deve ser exibido nos documentos clínicos e no perfil público do profissional.

---

## 15. Fora de escopo desta fase

- Lógica clínica (pacientes, prontuário, agenda, dietas)
- Geração de PDF
- Notificações
- Google Calendar
- IA assistiva
- Billing provider externo
- UI funcional nos portais
