# Frontend Shared Infrastructure — Design Spec

## Context

O backend MVP da Nutrometra (Fases 1-5) esta completo: 47 migrations, ~19.6k linhas Go, 19 endpoints de API REST cobrindo identity, tenancy, billing, scheduling, patients, clinical, diets, documents, exports, Google Calendar e IA assistiva.

Existem 3 apps Next.js como skeletons vazios (placeholder pages). Nenhuma escolha de frontend foi feita. Este spec define a infraestrutura compartilhada que as 3 apps usarao.

## Technology Stack

| Categoria | Escolha | Justificativa |
|-----------|---------|---------------|
| Framework | Next.js 14 (App Router) | Ja existente no projeto |
| Styling | Tailwind CSS | Utility-first, mobile-first nativo |
| Components | shadcn/ui (Radix primitives) | Copiavel, customizavel, acessivel |
| Data Fetching | TanStack Query (React Query) | Cache, invalidacao, retry, polling |
| Forms | React Hook Form + Zod | Performance em forms grandes, type-safe validation |
| Auth | next-auth (Auth.js) + Zitadel OIDC | Sessao via cookies HTTP-only, refresh automatico |
| Monorepo | Turborepo + pnpm workspaces | Build cache, tasks paralelas, dep resolution rapida |

## Monorepo Structure

```
/
├── services/api/                    # Go backend (existente, nao muda)
├── apps/
│   ├── web-professional/            # Next.js porta 3000
│   │   ├── src/
│   │   │   ├── app/                 # App Router pages
│   │   │   ├── components/          # Components especificos do portal profissional
│   │   │   └── lib/                 # Helpers especificos
│   │   ├── next.config.mjs
│   │   ├── tailwind.config.ts       # extends @nutrometra/ui/tailwind
│   │   ├── tsconfig.json
│   │   └── package.json
│   ├── web-patient/                 # Next.js porta 3001
│   │   └── (mesma estrutura, componentes especificos mobile-first)
│   └── backoffice/                  # Next.js porta 3002
│       └── (mesma estrutura, componentes especificos admin)
├── packages/
│   ├── ui/                          # @nutrometra/ui
│   │   ├── src/
│   │   │   ├── components/          # shadcn/ui components com tema Nutrometra
│   │   │   ├── layouts/             # AppShell, PageHeader, EmptyState, etc
│   │   │   └── tailwind/            # Preset compartilhado (cores, fontes, breakpoints)
│   │   ├── package.json
│   │   └── tsconfig.json
│   ├── api-client/                  # @nutrometra/api-client
│   │   ├── src/
│   │   │   ├── client.ts            # Fetch wrapper com auth + tenant headers
│   │   │   ├── types/               # Tipos TS espelhando domain types Go
│   │   │   ├── hooks/               # TanStack Query hooks por dominio
│   │   │   └── mutations/           # Mutation hooks com cache invalidation
│   │   ├── package.json
│   │   └── tsconfig.json
│   └── auth/                        # @nutrometra/auth
│       ├── src/
│       │   ├── config.ts            # NextAuth config com Zitadel provider
│       │   ├── middleware.ts         # Middleware re-exportavel para proteger rotas
│       │   └── hooks.ts             # useAuth, useTenant, usePermission, useEntitlements
│       ├── package.json
│       └── tsconfig.json
├── turbo.json
├── pnpm-workspace.yaml
├── package.json                     # Root workspace
├── .env.example                     # Atualizado com vars frontend
└── tsconfig.base.json               # Configuracao TS base compartilhada
```

## Package: @nutrometra/api-client

### Fetch wrapper (`client.ts`)

Um wrapper sobre `fetch` que:
- Injeta `Authorization: Bearer <token>` automaticamente (token da sessao next-auth)
- Injeta `X-Tenant-ID` header do contexto do tenant ativo
- Injeta `X-Request-ID` para correlacao com logs do backend
- Parseia respostas JSON e erros no formato `{code, message, request_id}` do backend
- Retry automatico para 5xx (via TanStack Query)
- Base URL configuravel via `NEXT_PUBLIC_API_URL`

### Tipos TypeScript (`types/`)

Tipos espelhando os domain types do Go backend. Organizados por dominio:

- `auth.ts` — User, Session, TokenPayload
- `tenancy.ts` — Tenant, TenantMembership
- `billing.ts` — Plan, Subscription, Entitlement, Invoice, Payment, FeatureOverride
- `professional.ts` — Professional, ProfessionalAddress, ServiceMode
- `scheduling.ts` — AvailabilityRule, ScheduleBlock, Appointment, TimeSlot
- `patient.ts` — Patient, PatientProfile, PatientInvite
- `clinical.ts` — Anamnesis, ProgressNote, ClinicalAttachment
- `bioimpedance.ts` — BodyMeasurement
- `catalog.ts` — Food, FoodGroup, HouseholdMeasure
- `diet.ts` — Diet, Meal, MealItem, Substitution, DietPublication
- `document.ts` — ClinicalDocument, DocumentVersion
- `export.ts` — ExportedFile
- `ai.ts` — AISuggestion, SuggestionType, SuggestionStatus
- `calendar.ts` — CalendarConnection, CalendarSyncLog
- `common.ts` — PaginatedResponse, ErrorResponse, UUID alias

### TanStack Query hooks (`hooks/`)

Um arquivo por dominio. Exemplo `hooks/patients.ts`:

```typescript
export function usePatients(params?: PatientListParams) {
  return useQuery({
    queryKey: ['patients', params],
    queryFn: () => api.get<PaginatedResponse<Patient>>('/patients', { params }),
  })
}

export function usePatient(patientId: string) {
  return useQuery({
    queryKey: ['patients', patientId],
    queryFn: () => api.get<Patient>(`/patients/${patientId}`),
    enabled: !!patientId,
  })
}
```

### Mutation hooks (`mutations/`)

```typescript
export function useCreatePatient() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreatePatientInput) => api.post<Patient>('/patients', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['patients'] })
    },
  })
}
```

## Package: @nutrometra/auth

### NextAuth Config (`config.ts`)

- Provider: Zitadel OIDC
- Session strategy: `jwt` (stateless, cookie-based)
- Callbacks:
  - `jwt`: extrair `tenant_id`, `roles`, `permissions` do token OIDC
  - `session`: expor tenant e roles na sessao client-side
- Token refresh automatico via `account.refresh_token`

### Middleware (`middleware.ts`)

- Exporta middleware configuravel por app
- Protege todas as rotas exceto `/auth/*` e assets estaticos
- Redireciona para login se sessao expirada
- Cada app importa e adapta: profissional valida role `nutritionist/owner`, paciente valida role `patient`, backoffice valida `backoffice_*`

### Hooks (`hooks.ts`)

- `useAuth()` — sessao, user, logout
- `useTenant()` — tenant ativo, switch tenant (para profissionais multi-tenant)
- `usePermission(code: string)` — boolean, consulta entitlements cacheados
- `useEntitlements()` — mapa completo de entitlements do tenant (cacheado via TanStack Query, TTL 5min alinhado com backend cache)

## Package: @nutrometra/ui

### Tailwind Preset

```typescript
// packages/ui/src/tailwind/preset.ts
export default {
  theme: {
    extend: {
      colors: {
        brand: { /* paleta Nutrometra */ },
        semantic: { success, warning, error, info },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
}
```

Cada app faz `presets: [require('@nutrometra/ui/tailwind')]` no seu tailwind.config.

### shadcn/ui Components

Inicializar com componentes essenciais:

**Inputs:** Button, Input, Textarea, Select, Checkbox, RadioGroup, Switch, DatePicker
**Layout:** Card, Separator, Sheet, Sidebar, Tabs, Accordion
**Feedback:** Alert, Badge, Toast, Skeleton, Progress
**Overlay:** Dialog, Dropdown, Popover, Tooltip, Command
**Data:** Table, Pagination
**Form:** Form (integracao shadcn + React Hook Form + Zod)

### Layout Components

- `AppShell` — sidebar + header + content area, responsive (sidebar collapsa em mobile)
- `PageHeader` — titulo + breadcrumbs + actions
- `EmptyState` — icone + mensagem + CTA (obrigatorio por CLAUDE.md)
- `LoadingState` — skeleton ou spinner contextual (obrigatorio)
- `ErrorState` — mensagem + retry (obrigatorio)
- `StaleIndicator` — badge visual quando dados estao potencialmente desatualizados (obrigatorio)

## Turborepo Configuration

### `turbo.json`

```json
{
  "$schema": "https://turbo.build/schema.json",
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": [".next/**", "dist/**"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    },
    "lint": {
      "dependsOn": ["^build"]
    },
    "test": {
      "dependsOn": ["^build"]
    }
  }
}
```

### `pnpm-workspace.yaml`

```yaml
packages:
  - "apps/*"
  - "packages/*"
```

## Environment Variables

### Novas vars frontend (adicionar ao `.env.example`)

```
# Frontend — API
NEXT_PUBLIC_API_URL=http://localhost:8081

# Frontend — Auth (Zitadel OIDC)
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=
ZITADEL_CLIENT_ID=
ZITADEL_CLIENT_SECRET=
ZITADEL_ISSUER=

# Frontend — Feature flags
NEXT_PUBLIC_ENABLE_AI=true
NEXT_PUBLIC_ENABLE_GOOGLE_CALENDAR=true
```

Valores sensiveis (`NEXTAUTH_SECRET`, `ZITADEL_CLIENT_SECRET`) ficam SOMENTE no `.env` (gitignored). O `.env.example` tem campos vazios.

## Migration das Apps Existentes

As 3 apps atuais (`web-professional/`, `web-patient/`, `backoffice/`) serao:
1. Movidas para `apps/`
2. Atualizadas para usar pnpm e importar dos packages compartilhados
3. `package.json` atualizado com dependencias: `@nutrometra/ui`, `@nutrometra/api-client`, `@nutrometra/auth`, `tailwindcss`, `next-auth`, `@tanstack/react-query`, `react-hook-form`, `zod`

## Verificacao

1. `pnpm install` — resolve todas as dependencias
2. `pnpm turbo build` — compila os 3 packages e as 3 apps
3. `pnpm turbo dev` — inicia as 3 apps em paralelo (3000, 3001, 3002)
4. Cada app mostra pagina com layout basico (AppShell + auth redirect)
5. Tipos TS compilam sem erros
6. Tailwind gera classes corretamente com tema compartilhado

## Fora de escopo

- Implementacao de telas/features (sera feito nos proximos specs por portal)
- API OpenAPI spec (pode ser gerado depois)
- Testes E2E (serao adicionados com cada portal)
- Internacionalizacao (portugues fixo por ora)
- PWA/offline (futuro)
