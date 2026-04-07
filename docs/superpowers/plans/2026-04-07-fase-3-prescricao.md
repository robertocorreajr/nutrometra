# Fase 3 — Prescrição e Publicação

## Contexto

As Fases 1 (fundação) e 2 (operação clínica) estão completas e merged em `main`. A Fase 3 implementa os três módulos obrigatórios restantes para completar o ciclo clínico: **dietas por refeições**, **documentos clínicos** e **geração de PDF**. Sem estes módulos, o profissional não consegue prescrever, o paciente não consegue visualizar prescrições, e nenhum documento pode ser exportado — funcionalidades centrais do produto.

Referência: `docs/05-delivery/mvp-phases.md` §5 — "Fase 3 — Prescrição e publicação obrigatória".

**Implemented in:** PR #5 (`feat/fase-3-prescricao`)

---

## Ajustes aplicados durante validação

O plano original foi validado contra o codebase atual. **10 ajustes** foram identificados e aplicados:

### A1. Referência incorreta ao padrão de supersede
O plano original referenciava `bioimpedance/repository/postgres.go:141` como padrão de supersede. Na realidade, esse código retorna `ErrAlreadyPublished` em conflito de unique constraint — **não** implementa supersede. O `SupersedePublications` do módulo diet é **código novo**.

### A2. Bug de URL param em rotas de paciente
As rotas em `main.go:226` definem `/{id}` para pacientes, mas handlers como `clinical/handler.go:110` leem `chi.URLParam(r, "patient_id")`. Isso é um bug latente. **Correção**: renomear a rota para `/{patient_id}` e verificar consistência nos handlers existentes.

### A3. Padrão de transação no usecase
Nenhum usecase existente recebe `pool`. Transações são gerenciadas no handler via `db.RunInTx`. Porém, `PublishDiet` e `CreateNewVersion` são operações multi-step complexas demais para o handler orquestrar. **Decisão**: diet e document usecases recebem `pool` — desvio justificado e documentado como evolução do padrão.

### A4. Biblioteca PDF
`jung-kurt/gofpdf` está archived. Usar `github.com/go-pdf/fpdf` (fork mantido). É a única dependência nova.

### A5. Shutdown do worker
O worker deve parar **antes** do server no shutdown sequence: `bgWorker.Stop()` → `srv.Shutdown()`.

### A6. Permissão do paciente para documentos
O role `patient` não tem `clinical:read` (tem `self:read`, `diet:read`, `appointments:read`). Adicionar `clinical:read` ao patient na migration 000037 para que pacientes vejam documentos publicados.

### A7. Índice parcial em diet_publications
Adicionar `CREATE UNIQUE INDEX idx_diet_pub_active ON diet_publications(diet_id) WHERE status = 'active'` para prevenir race conditions sem necessidade de `SELECT FOR UPDATE`.

### A8. Import aliases no wiring
Explicitar aliases: `dietrepo`, `dietuc`, `docrepo`, `docuc`, `exportrepo`, `exportuc`.

### A9. Construtor do export usecase
O exportUC com 5 dependências é o mais complexo do codebase. Justificativa: ele orquestra entitlement check + criação de registro + enfileiramento de job + geração de PDF. Manter assim.

### A10. Status em document_publications
Documentos clínicos **não** suportam supersede — cada versão é independentemente publicada. Sem coluna `status` em `document_publications`. Decisão explícita.

---

## Escopo

| Módulo | Descrição |
|--------|-----------|
| **diet** | Plano alimentar com refeições, itens, substituições, versionamento e publicação ao paciente |
| **document** | Documentos clínicos (solicitação de exames, prescrições, cartas), versionamento, finalização e publicação |
| **export** | Geração assíncrona de PDF, tracking de status, download, histórico de exportações |
| **platform/worker** | Worker goroutine-based para jobs assíncronos (PDF) |

---

## Step 0 — Bugfix de URL param (pré-requisito)

**Arquivo**: `services/api/cmd/api/main.go:226`

Renomear `/{id}` para `/{patient_id}` no bloco de rotas de pacientes. Verificar e corrigir handlers existentes que leem `chi.URLParam(r, "id")` dentro desse bloco (patient, clinical, bioimpedance).

---

## Step 1 — Migrations (000032–000037)

Branch: `feat/fase-3-prescricao`

### 000032_create_diets.up.sql
```sql
-- diets: plano alimentar do paciente
CREATE TABLE diets (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id),
    patient_id          UUID NOT NULL REFERENCES patients(id),
    professional_id     UUID NOT NULL REFERENCES professionals(id),
    title               TEXT NOT NULL,
    objective           TEXT,
    status              TEXT NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'published', 'archived')),
    version_number      INT NOT NULL DEFAULT 1,
    previous_version_id UUID REFERENCES diets(id),
    published_at        TIMESTAMPTZ,
    valid_from          DATE,
    valid_until         DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_diets_tenant ON diets(tenant_id);
CREATE INDEX idx_diets_patient ON diets(tenant_id, patient_id);
CREATE INDEX idx_diets_professional ON diets(tenant_id, professional_id);
```

### 000033_create_diet_meals_items.up.sql
```sql
-- diet_meals: refeições dentro de uma dieta (café da manhã, almoço, etc.)
CREATE TABLE diet_meals (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    diet_id     UUID NOT NULL REFERENCES diets(id) ON DELETE CASCADE,
    meal_name   TEXT NOT NULL,
    meal_order  INT NOT NULL,
    notes       TEXT,
    UNIQUE (diet_id, meal_order)
);
CREATE INDEX idx_diet_meals_diet ON diet_meals(diet_id);

-- diet_meal_items: itens alimentares dentro de cada refeição
CREATE TABLE diet_meal_items (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id),
    diet_meal_id         UUID NOT NULL REFERENCES diet_meals(id) ON DELETE CASCADE,
    food_item_id         UUID NOT NULL REFERENCES food_items(id),
    quantity_value       NUMERIC(8,2) NOT NULL,
    quantity_unit        TEXT NOT NULL,
    household_measure_id UUID REFERENCES household_measures(id),
    amount_description   TEXT,
    preparation_notes    TEXT,
    sort_order           INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_diet_meal_items_meal ON diet_meal_items(diet_meal_id);

-- diet_substitutions: substituições possíveis para um item
CREATE TABLE diet_substitutions (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants(id),
    diet_meal_item_id       UUID NOT NULL REFERENCES diet_meal_items(id) ON DELETE CASCADE,
    substitute_food_item_id UUID NOT NULL REFERENCES food_items(id),
    quantity_value          NUMERIC(8,2),
    quantity_unit           TEXT,
    notes                   TEXT,
    sort_order              INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_diet_subs_item ON diet_substitutions(diet_meal_item_id);
```

### 000034_create_diet_publications.up.sql
```sql
CREATE TABLE diet_publications (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id),
    diet_id              UUID NOT NULL REFERENCES diets(id),
    patient_id           UUID NOT NULL REFERENCES patients(id),
    published_by_user_id UUID NOT NULL REFERENCES users(id),
    published_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status               TEXT NOT NULL DEFAULT 'active'
                         CHECK (status IN ('active', 'superseded'))
);
CREATE INDEX idx_diet_pub_patient ON diet_publications(tenant_id, patient_id);
-- A7: índice parcial único para prevenir race conditions
CREATE UNIQUE INDEX idx_diet_pub_active ON diet_publications(diet_id) WHERE status = 'active';
```

### 000035_create_clinical_documents.up.sql
```sql
CREATE TABLE clinical_documents (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id),
    patient_id          UUID NOT NULL REFERENCES patients(id),
    professional_id     UUID NOT NULL REFERENCES professionals(id),
    appointment_id      UUID REFERENCES appointments(id),
    document_type       TEXT NOT NULL
                        CHECK (document_type IN ('exam_request', 'prescription', 'letter', 'other')),
    title               TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'finalized', 'published')),
    content_json        JSONB NOT NULL DEFAULT '{}',
    version_number      INT NOT NULL DEFAULT 1,
    previous_version_id UUID REFERENCES clinical_documents(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_clin_docs_patient ON clinical_documents(tenant_id, patient_id);

-- A10: sem coluna status em document_publications (cada versão é independentemente publicada)
CREATE TABLE document_publications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants(id),
    clinical_document_id    UUID NOT NULL REFERENCES clinical_documents(id),
    patient_id              UUID NOT NULL REFERENCES patients(id),
    published_by_user_id    UUID NOT NULL REFERENCES users(id),
    published_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_doc_pub_patient ON document_publications(tenant_id, patient_id);
```

### 000036_create_exported_files.up.sql
```sql
CREATE TABLE exported_files (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id),
    related_entity_type  TEXT NOT NULL,
    related_entity_id    UUID NOT NULL,
    export_type          TEXT NOT NULL DEFAULT 'pdf'
                         CHECK (export_type IN ('pdf', 'print_job')),
    file_key             TEXT,
    status               TEXT NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    requested_by_user_id UUID NOT NULL REFERENCES users(id),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at         TIMESTAMPTZ,
    failure_reason       TEXT
);
CREATE INDEX idx_exported_files_entity ON exported_files(related_entity_type, related_entity_id);
CREATE INDEX idx_exported_files_status ON exported_files(tenant_id, status);
```

### 000037_seed_phase3_permissions.up.sql

Novas permissões (IDs sequenciais após `..0018`):

```sql
INSERT INTO permissions (id, code, application_scope, description) VALUES
  ('10000000-0000-0000-0000-000000000030', 'diet:write',      'tenant', 'Criar e editar dietas'),
  ('10000000-0000-0000-0000-000000000031', 'diet:publish',    'tenant', 'Publicar dietas ao paciente'),
  ('10000000-0000-0000-0000-000000000032', 'document:write',  'tenant', 'Criar e editar documentos clínicos'),
  ('10000000-0000-0000-0000-000000000033', 'document:publish','tenant', 'Publicar documentos ao paciente'),
  ('10000000-0000-0000-0000-000000000034', 'export:pdf',      'tenant', 'Exportar PDF')
ON CONFLICT (id) DO NOTHING;

-- owner: recebe automaticamente (já usa SELECT WHERE scope='tenant')
-- nutritionist: precisa de diet:write, diet:publish, document:write, document:publish, export:pdf
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000002', id
FROM permissions WHERE code IN ('diet:write','diet:publish','document:write','document:publish','export:pdf')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- A6: patient precisa de clinical:read para ver documentos publicados
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000004', id
FROM permissions WHERE code = 'clinical:read'
ON CONFLICT (role_id, permission_id) DO NOTHING;
```

Down migrations: `DROP TABLE` em ordem inversa; `DELETE FROM permissions WHERE code IN (...)`.

---

## Step 2 — Módulo Diet

**Arquivos:**
- `services/api/internal/diet/domain/diet.go`
- `services/api/internal/diet/domain/diet_test.go`
- `services/api/internal/diet/repository/postgres.go`
- `services/api/internal/diet/usecase/diet.go`
- `services/api/internal/diet/handler.go`

### Domain (`diet.go`)

Tipos: `Diet`, `DietMeal`, `DietMealItem`, `DietSubstitution`, `DietPublication`

Status machine:
```
draft ──► published ──► archived
  │                        ▲
  └── (delete ok)          │
       published ──────────┘ (new-version cria draft novo, arquiva o publicado)
```

Sentinel errors: `ErrNotFound`, `ErrDietNotDraft`, `ErrDietEmpty` (sem refeições/itens), `ErrAlreadyPublished`, `ErrNotPublished`

Validação:
- `Diet.Validate()`: título obrigatório
- `DietMeal.Validate()`: meal_name e meal_order > 0
- `DietMealItem.Validate()`: quantity_value > 0, quantity_unit obrigatório
- `Diet.HasContent()`: verifica se tem pelo menos 1 meal com 1 item

### Repository (`postgres.go`)

Usa `*pgxpool.Pool` como campo principal. Métodos que precisam de transação recebem `pgx.Tx` como parâmetro (mesmo padrão de `db.RunInTx`). Interface `querier` satisfeita tanto por `*pgxpool.Pool` quanto por `pgx.Tx`.

Full graph loading via 4 queries separadas (diet, meals, items, substitutions) — sem JOINs complexos.

Métodos:
- `CreateDiet(ctx, tx, *Diet)` / `GetByID(ctx, tenantID, id)` (carrega grafo completo: meals→items→subs)
- `UpdateDiet(ctx, *Diet)` / `DeleteDiet(ctx, tenantID, id)` / `ListByPatient(ctx, tenantID, patientID)`
- `CreateMeal(ctx, tx, *DietMeal)` / `UpdateMeal` / `DeleteMeal`
- `CreateMealItem(ctx, tx, *DietMealItem)` / `UpdateMealItem` / `DeleteMealItem`
- `CreateSubstitution(ctx, tx, *DietSubstitution)` / `DeleteSubstitution`
- `CreatePublication(ctx, tx, *DietPublication)`
- `SupersedePublications(ctx, tx, tenantID, patientID)` — marca publicações ativas como `superseded` (A1: código novo)
- Métodos transacionais para deep-copy: `GetDietForVersion(tx)`, `ArchiveDietTx(tx)`, `CreateDietTx(tx)`, `CreateMealTx(tx)`, `CreateMealItemTx(tx)`, `CreateSubstitutionTx(tx)`

### Usecase (`diet.go`)

Interface do repositório definida localmente. Pool recebido no construtor para `db.RunInTx` (A3: evolução do padrão).

Construtor: `New(repo DietRepository, pool *pgxpool.Pool) *Usecase`

Lógica crítica:

**PublishDiet**: (dentro de transação)
1. Carrega dieta, valida status = `draft`
2. Valida que existe pelo menos 1 meal com 1 item
3. Atualiza status → `published`, set `published_at`
4. Supersede publicações ativas do mesmo paciente
5. Cria novo registro em `diet_publications`

**CreateNewVersion**: (dentro de transação)
1. Carrega dieta completa (meals + items + subs) dentro da transação
2. Valida status = `published` ou `archived`
3. Cria nova dieta com `previous_version_id`, `version_number + 1`, status `draft`
4. Deep-copy de meals → items → substitutions com novos UUIDs
5. Arquiva a versão anterior se estava `published`

**ArchiveDiet**: Valida status = `published`, muda para `archived`

### Handler (`handler.go`)

Segue o padrão exato dos outros handlers (extract tenantID/userID do ctx, decode body, call usecase, audit, render). Construtor: `NewHandler(uc, pool, auditSvc)`, lê `chi.URLParam(r, "patient_id")`.

### Endpoints (16 rotas)

| Método | Path | Permissão |
|--------|------|-----------|
| POST | `/patients/{patient_id}/diets` | `diet:write` |
| GET | `/patients/{patient_id}/diets` | `diet:read` |
| GET | `/diets/{id}` | `diet:read` |
| PUT | `/diets/{id}` | `diet:write` |
| DELETE | `/diets/{id}` | `diet:write` |
| POST | `/diets/{id}/publish` | `diet:publish` |
| POST | `/diets/{id}/archive` | `diet:write` |
| POST | `/diets/{id}/new-version` | `diet:write` |
| POST | `/diets/{id}/meals` | `diet:write` |
| PUT | `/diets/{id}/meals/{meal_id}` | `diet:write` |
| DELETE | `/diets/{id}/meals/{meal_id}` | `diet:write` |
| POST | `/diets/{id}/meals/{meal_id}/items` | `diet:write` |
| PUT | `/diet-items/{item_id}` | `diet:write` |
| DELETE | `/diet-items/{item_id}` | `diet:write` |
| POST | `/diet-items/{item_id}/substitutions` | `diet:write` |
| DELETE | `/diet-substitutions/{sub_id}` | `diet:write` |

---

## Step 3 — Módulo Document

**Arquivos:**
- `services/api/internal/document/domain/document.go`
- `services/api/internal/document/domain/document_test.go`
- `services/api/internal/document/repository/postgres.go`
- `services/api/internal/document/usecase/document.go`
- `services/api/internal/document/handler.go`

### Domain

Tipos: `ClinicalDocument`, `DocumentPublication`

Document types: `exam_request`, `prescription`, `letter`, `other`

Status machine:
```
draft ──► finalized ──► published
  │                        ▲
  └── (delete ok)          │ (new-version cria draft)
```

Diferença da diet: documentos têm step intermediário `finalized` (conteúdo travado, irreversível) antes de publicar. Segue o fluxo descrito em `docs/02-domain/clinical-workflows.md` §8.

Sentinel errors: `ErrNotFound`, `ErrNotDraft`, `ErrNotFinalized`, `ErrAlreadyFinalized`

### Repository

- `Create`, `GetByID`, `Update` (draft only), `ListByPatient`
- `Finalize(ctx, tenantID, id)` — set status finalized, updated_at
- `Publish(ctx, tx, tenantID, id)` — set status published
- `CreatePublication(ctx, tx, *DocumentPublication)`
- `ListVersions(ctx, tenantID, id)` — recursive CTE para seguir cadeia `previous_version_id`

A10: `document_publications` sem coluna `status` — cada versão é independentemente publicada, sem supersede.

### Usecase

Construtor: `New(repo DocumentRepository, pool *pgxpool.Pool) *Usecase` (A3: recebe pool)

**FinalizeDocument**: Valida draft → finalized (irreversível)
**PublishDocument**: Valida finalized → published + cria publication (transacional)
**CreateNewVersion**: Deep-copy com previous_version_id, version_number + 1

### Endpoints (8 rotas)

| Método | Path | Permissão |
|--------|------|-----------|
| POST | `/patients/{patient_id}/documents` | `document:write` |
| GET | `/patients/{patient_id}/documents` | `clinical:read` |
| GET | `/documents/{id}` | `clinical:read` |
| PUT | `/documents/{id}` | `document:write` |
| POST | `/documents/{id}/finalize` | `document:write` |
| POST | `/documents/{id}/publish` | `document:publish` |
| POST | `/documents/{id}/new-version` | `document:write` |
| GET | `/documents/{id}/versions` | `clinical:read` |

---

## Step 4 — Módulo Export (PDF)

**Arquivos:**
- `services/api/internal/export/domain/export.go`
- `services/api/internal/export/repository/postgres.go`
- `services/api/internal/export/usecase/export.go`
- `services/api/internal/export/handler.go`

### Domain

`ExportedFile` com status: `pending` → `processing` → `completed`/`failed`

### Usecase — Fluxo assíncrono

1. Handler recebe `POST /exports` com `{entity_type, entity_id, export_type}`
2. Usecase checa entitlement `pdf:export` via `EntitlementAdapter`
3. Cria registro `exported_files` com status `pending`
4. Enfileira job no worker
5. Retorna **202 Accepted** com ID do export
6. Worker: `processing` → gera PDF → `completed` com `file_key` / ou `failed`
7. Cliente faz polling em `GET /exports/{id}`
8. Download via `GET /exports/{id}/download`

A9: exportUC com 5 dependências (exportRepo, pool, entitlementAdapter, pdfGen, bgWorker) — mais complexo do codebase, justificado pela orquestração.

### Endpoints (4 rotas)

| Método | Path | Permissão |
|--------|------|-----------|
| POST | `/exports` | `export:pdf` |
| GET | `/exports` | `export:pdf` |
| GET | `/exports/{id}` | `export:pdf` |
| GET | `/exports/{id}/download` | `export:pdf` |

---

## Step 5 — Infraestrutura de Worker e PDF

### platform/worker (`services/api/internal/platform/worker/worker.go`)

Worker simples baseado em goroutine + channel, executando dentro do processo da API:

```go
type Job struct { ID uuid.UUID; Type string; Payload any }
type Handler func(ctx context.Context, job Job) error
type Worker struct { handlers map[string]Handler; queue chan Job }
```

- Buffer de 100 jobs
- Criado em `main.go`, registra handler de PDF, inicia como goroutine
- `Start(ctx)` e `Stop()` com drain via sync.WaitGroup
- A5: Worker shutdown **antes** do server: `bgWorker.Stop()` → `srv.Shutdown()`
- Evolução futura: trocar channel por Redis LIST (LPUSH/BRPOP)

### Geração de PDF (`services/api/internal/platform/pdfgen/pdfgen.go`)

**Abordagem MVP**: Usar biblioteca Go pura (`go-pdf/fpdf`, A4: fork mantido) para gerar PDF diretamente sem dependência externa. Não precisa de Chrome/wkhtmltopdf.

**Alternativa**: Se a qualidade visual for insuficiente, migrar para `chromedp` (headless Chrome) com HTML templates. O design permite essa troca transparente — a interface `Generator` isola a implementação.

```go
type Generator interface {
    GenerateDietPDF(ctx context.Context, data DietPDFData) ([]byte, error)
    GenerateDocumentPDF(ctx context.Context, data DocumentPDFData) ([]byte, error)
}
```

Implementação `FPDFGenerator` gera PDFs com:
- Header com título e informações do paciente/profissional
- Seções de refeições com itens e substituições (diet)
- Seções key-value do content map (document)
- Footer com data de geração

**Storage**: Arquivos salvos em filesystem local com path `exports/{tenant_id}/{export_id}.pdf`. O campo `file_key` guarda esse path. Migração para object storage (GCS) é transparente — basta trocar a implementação de storage.

---

## Step 6 — Wiring em main.go

**Arquivo**: `services/api/cmd/api/main.go`

Imports (A8: aliases explícitos):
```go
"nutrometra/api/internal/diet"
dietrepo "nutrometra/api/internal/diet/repository"
dietuc "nutrometra/api/internal/diet/usecase"
"nutrometra/api/internal/document"
docrepo "nutrometra/api/internal/document/repository"
docuc "nutrometra/api/internal/document/usecase"
"nutrometra/api/internal/export"
exportrepo "nutrometra/api/internal/export/repository"
exportuc "nutrometra/api/internal/export/usecase"
"nutrometra/api/internal/platform/worker"
"nutrometra/api/internal/platform/pdfgen"
```

Adicionar na seção de repositories:
```go
dietRepo   := dietrepo.New(pool)
docRepo    := docrepo.New(pool)
exportRepo := exportrepo.New(pool)
```

Adicionar na seção de services:
```go
pdfGen    := pdfgen.New()
bgWorker  := worker.New(100)
dietUC    := dietuc.New(dietRepo, pool)
docUC     := docuc.New(docRepo, pool)
exportUC  := exportuc.New(exportRepo, pool, entitlementAdapter, pdfGen, bgWorker)
```

Adicionar na seção de handlers:
```go
dietHandler   := diet.NewHandler(dietUC, pool, auditSvc)
docHandler    := document.NewHandler(docUC, pool, auditSvc)
exportHandler := export.NewHandler(exportUC, pool, auditSvc)
```

Lifecycle (A5):
```go
bgWorker.Start(ctx)
// ... server start em goroutine ...
<-quit
bgWorker.Stop()
srv.Shutdown(shutdownCtx)
```

Rotas: diets e documents aninhados sob `/patients/{patient_id}/...`, resources standalone sob `/diets/{id}/...`, `/documents/{id}/...`, `/exports/...`.

---

## Step 7 — Testes

### Testes unitários (domain)
- `diet/domain/diet_test.go` — Validate() para Diet, DietMeal, DietMealItem; HasContent(); transições de status
- `document/domain/document_test.go` — Validate() e transições draft→finalized→published

### Testes de usecase (mock do repo)
- PublishDiet: valida transição, rejeita dieta vazia, supersede publicações anteriores
- CreateNewVersion: valida deep-copy, incremento de version_number
- FinalizeDocument: irreversível, rejeita re-finalização
- PublishDocument: requer finalized, cria publication

---

## Auditoria

Todas as operações mutáveis registram audit seguindo o padrão existente:

| Entity Type | Actions |
|-------------|---------|
| `diet` | `created`, `updated`, `deleted`, `published`, `archived`, `version_created` |
| `diet_meal` | `added`, `updated`, `deleted` |
| `diet_meal_item` | `added`, `updated`, `deleted` |
| `diet_substitution` | `added`, `deleted` |
| `clinical_document` | `created`, `updated`, `finalized`, `published`, `version_created` |
| `exported_file` | `requested`, `completed`, `failed` |

---

## Ordem de implementação

1. **Migrations** (000032–000037) — sem dependência
2. **Diet domain + tests** — sem dependência
3. **Diet repository** — depende de 1
4. **Diet usecase** — depende de 2, 3
5. **Diet handler** — depende de 4
6. **Document domain + tests** — paralelo a 2-5
7. **Document repository** — depende de 1
8. **Document usecase + handler** — depende de 6, 7
9. **Export domain + repository** — depende de 1
10. **Worker + PDFgen** — sem dependência de módulos
11. **Export usecase + handler** — depende de 9, 10
12. **Wiring em main.go** — depende de 5, 8, 11
13. **Verificação end-to-end** — depende de 12

---

## Mapa de arquivos

### Diet
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/diet/domain/diet.go` | tipos Diet, DietMeal, DietMealItem, DietSubstitution, DietPublication |
| `services/api/internal/diet/domain/diet_test.go` | testes de validação e HasContent |
| `services/api/internal/diet/repository/postgres.go` | CRUD + full graph loading (4 queries) + SupersedePublications |
| `services/api/internal/diet/usecase/diet.go` | PublishDiet e CreateNewVersion transacionais |
| `services/api/internal/diet/handler.go` | 16 endpoints |

### Document
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/document/domain/document.go` | tipos ClinicalDocument, DocumentPublication |
| `services/api/internal/document/domain/document_test.go` | testes de validação |
| `services/api/internal/document/repository/postgres.go` | CRUD + ListVersions (recursive CTE) |
| `services/api/internal/document/usecase/document.go` | Finalize irreversível, Publish e CreateNewVersion transacionais |
| `services/api/internal/document/handler.go` | 8 endpoints |

### Export
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/export/domain/export.go` | tipo ExportedFile com status machine |
| `services/api/internal/export/repository/postgres.go` | CRUD + UpdateStatus |
| `services/api/internal/export/usecase/export.go` | RequestExport com entitlement check + job dispatch |
| `services/api/internal/export/handler.go` | 4 endpoints (POST retorna 202 Accepted) |

### Platform
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/platform/worker/worker.go` | goroutine + channel (buffer 100), Start/Stop com drain |
| `services/api/internal/platform/pdfgen/pdfgen.go` | interface Generator + implementação FPDFGenerator |

### Migrations (000032–000037)
| Arquivo | Conteúdo |
|---------|----------|
| `000032_create_diets` | tabela diets com versionamento e status machine |
| `000033_create_diet_meals_items` | diet_meals, diet_meal_items, diet_substitutions |
| `000034_create_diet_publications` | publicações com índice parcial único |
| `000035_create_clinical_documents` | clinical_documents + document_publications |
| `000036_create_exported_files` | tracking de exports assíncronos |
| `000037_seed_phase3_permissions` | 5 permissões + clinical:read para patient |

---

## Arquivos criados (~30 novos, 1 modificado)

| Categoria | Qtd | Detalhes |
|-----------|-----|---------|
| Migrations (novos) | 12 | 000032–000037 (.up.sql + .down.sql) |
| Diet module (novos) | 5 | domain, domain_test, repository, usecase, handler |
| Document module (novos) | 5 | domain, domain_test, repository, usecase, handler |
| Export module (novos) | 4 | domain, repository, usecase, handler |
| Platform (novos) | 2 | worker/worker.go, pdfgen/pdfgen.go |
| Modificados | 1 | cmd/api/main.go (wiring + bugfix route param) |
| **Total** | **29** | 28 novos + 1 modificado |

---

## Status machines

### Diet
```
draft → published → archived
  │                    ▲
  └── (delete ok)      │
       published ──────┘ (new-version cria draft, arquiva publicado)
```

### Document
```
draft → finalized → published
  │ (delete ok)
  │ finalized é IRREVERSÍVEL
  │ published permite new-version (cria novo draft)
```

### Export
```
pending → processing → completed (com file_key)
                     → failed (com failure_reason)
```

---

## Padrões reutilizados do codebase existente

- `db.RunInTx()` → `services/api/internal/platform/db/db.go:32`
- `audit.NewEntry()` com options → `services/api/internal/platform/audit/audit.go:48`
- `server.RenderJSON()` / `server.RenderError()` → `services/api/internal/platform/server/server.go:86-102`
- `identitydomain.TenantIDFromContext()` / `UserIDFromContext()` → `services/api/internal/identity/domain/user.go`
- `rbac.RequirePermission()` → `services/api/internal/rbac/middleware.go:21`
- `EntitlementAdapter.CheckEntitlement()` → `services/api/internal/billing/usecase/adapter.go:22`
- Handler pattern: `NewHandler(uc, pool, auditSvc)` → todos os handlers da Fase 2

---

## Verificação

- [x] `go build ./...` — compila sem erros
- [x] `go vet ./...` — sem warnings
- [x] `go test ./...` — todos os testes passam
- [x] Docker Compose: API sobe, worker started, health check OK
- [x] Migrations aplicam em banco limpo
- [x] Endpoints retornam 401 sem token
- [x] Fluxo: criar dieta → refeições → itens → substituições → publicar → nova versão
- [x] Fluxo: criar documento → finalizar → publicar → nova versão
- [x] Fluxo: solicitar export PDF → poll status → download

---

## Riscos e mitigações

| Risco | Mitigação |
|-------|-----------|
| Bugfix de URL param quebra handlers existentes | Verificar todos os `chi.URLParam` nos handlers de patient, clinical e bioimpedance |
| Deep-copy de dieta falha parcialmente | Transação explícita com rollback |
| PDFs perdidos em restart | `exported_files` é source of truth; jobs `pending`/`processing` podem ser re-enfileirados no startup |
| Filesystem local é efêmero em container | Documentado como limitação; `file_key` é abstrato — migração para GCS é transparente |
| Publicação concorrente | Índice parcial único `idx_diet_pub_active` (A7) |
| Dependência externa de PDF | MVP usa biblioteca Go pura (sem Chrome/wkhtmltopdf) |
