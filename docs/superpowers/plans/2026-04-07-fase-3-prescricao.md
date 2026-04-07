# Fase 3 — Prescrição e Publicação: Plano de Implementação

**Goal:** Completar o ciclo clínico com dietas por refeições, documentos clínicos e geração de PDF. Sem estes módulos, o profissional não prescreve, o paciente não visualiza prescrições e nenhum documento pode ser exportado.

**Architecture:** Segue os padrões das Fases 1-2 com uma evolução: diet e document usecases recebem `pool` para transações multi-step (PublishDiet, CreateNewVersion). Worker goroutine-based para jobs assíncronos (PDF). Interface `PDFGenerator` com implementação via `go-pdf/fpdf`.

**Implemented in:** PR #5 (`feat/fase-3-prescricao`)

---

## Ajustes aplicados durante validação

1. **A1**: Padrão de supersede é código novo (bioimpedance não implementava supersede)
2. **A2**: Bugfix de URL param `/{id}` → `/{patient_id}` nas rotas de paciente
3. **A3**: Usecases recebem `pool` para transações complexas (evolução do padrão)
4. **A4**: Biblioteca PDF: `go-pdf/fpdf` (fork mantido, `jung-kurt/gofpdf` archived)
5. **A5**: Worker shutdown antes do server no graceful shutdown
6. **A6**: Permissão `clinical:read` adicionada ao role patient
7. **A7**: Índice parcial único em `diet_publications` para prevenir race conditions
8. **A8**: Import aliases explícitos (dietrepo, dietuc, docrepo, docuc, etc.)
9. **A9**: Export usecase com 5 deps (mais complexo do codebase, justificado)
10. **A10**: `document_publications` sem supersede (cada versão publicada independentemente)

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

## Endpoints implementados

### Diet (16 rotas)
- `POST/GET /patients/{patient_id}/diets` — criar e listar
- `GET/PUT/DELETE /diets/{id}` — CRUD standalone
- `POST /diets/{id}/publish` — publicar (supersede ativas)
- `POST /diets/{id}/archive` — arquivar
- `POST /diets/{id}/new-version` — nova versão (deep copy)
- `POST /diets/{id}/meals` — adicionar refeição
- `PUT/DELETE /diets/{id}/meals/{meal_id}` — editar/remover refeição
- `POST /diets/{id}/meals/{meal_id}/items` — adicionar item
- `PUT/DELETE /diet-items/{item_id}` — editar/remover item
- `POST /diet-items/{item_id}/substitutions` — adicionar substituição
- `DELETE /diet-substitutions/{sub_id}` — remover substituição

### Document (8 rotas)
- `POST/GET /patients/{patient_id}/documents` — criar e listar
- `GET/PUT /documents/{id}` — detalhe e edição
- `POST /documents/{id}/finalize` — finalização irreversível
- `POST /documents/{id}/publish` — publicar (requer finalizado)
- `POST /documents/{id}/new-version` — nova versão
- `GET /documents/{id}/versions` — histórico de versões

### Export (4 rotas)
- `POST /exports` — solicitar export (202 Accepted)
- `GET /exports` — listar exports
- `GET /exports/{id}` — status do export
- `GET /exports/{id}/download` — download do PDF

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

## Verificação

- [x] `go build ./...` compila sem erros
- [x] `go vet ./...` sem warnings
- [x] `go test ./...` — todos os testes passam
- [x] Docker Compose: API sobe, worker started, health check OK
- [x] Migrations aplicam em banco limpo
- [x] Endpoints retornam 401 sem token
