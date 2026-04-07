# Fase 2 — Operação Clínica: Plano de Implementação

**Goal:** Camada de operação clínica completa — profissionais, agenda, pacientes, prontuário, bioimpedância e catálogo alimentar. Todos os módulos com CRUD, RBAC, auditoria e isolamento multi-tenant.

**Architecture:** Segue os padrões da Fase 1. Cada módulo segue a estrutura `domain → repository → usecase → handler`. Handlers com `NewHandler(uc, pool, auditSvc)`. RBAC por rota via `rbac.RequirePermission()`. Audit em todas as operações mutáveis.

**Implemented in:** PR #3 (`120f424`), PR #4 (`b9b76cd`)

---

## Mapa de arquivos

### Professionals
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/professional/domain/professional.go` | tipos Professional, Address, ServiceMode |
| `services/api/internal/professional/domain/professional_test.go` | testes de validação |
| `services/api/internal/professional/repository/postgres.go` | CRUD professionals, addresses, service_modes |
| `services/api/internal/professional/usecase/professional.go` | lógica de negócio com entitlement check |
| `services/api/internal/professional/handler.go` | endpoints /professionals |

### Scheduling
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/scheduling/domain/scheduling.go` | tipos AvailabilityRule, ScheduleBlock, Appointment |
| `services/api/internal/scheduling/repository/postgres.go` | CRUD agenda, bloqueios, agendamentos |
| `services/api/internal/scheduling/usecase/scheduling.go` | lógica de disponibilidade e lifecycle de appointments |
| `services/api/internal/scheduling/handler.go` | endpoints /appointments, /availability, /blocks, /slots |

### Patients
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/patient/domain/patient.go` | tipos Patient, PatientInvite, PatientProfile |
| `services/api/internal/patient/repository/postgres.go` | CRUD patients, invites, profiles, access_links |
| `services/api/internal/patient/usecase/patient.go` | lógica com entitlement check, convites |
| `services/api/internal/patient/handler.go` | endpoints /patients, /invites |

### Clinical
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/clinical/domain/clinical.go` | tipos Anamnesis, ProgressNote, ClinicalAttachment |
| `services/api/internal/clinical/repository/postgres.go` | CRUD anamneses, notas, anexos |
| `services/api/internal/clinical/usecase/clinical.go` | lógica com finalização irreversível |
| `services/api/internal/clinical/handler.go` | endpoints /anamneses, /notes, /attachments |

### Bioimpedance
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/bioimpedance/domain/measurement.go` | tipos BodyMeasurement, MeasurementPublication |
| `services/api/internal/bioimpedance/repository/postgres.go` | CRUD medições, publicações |
| `services/api/internal/bioimpedance/usecase/bioimpedance.go` | lógica de publicação |
| `services/api/internal/bioimpedance/handler.go` | endpoints /measurements |

### Food Catalog
| Arquivo | Responsabilidade |
|---------|-----------------|
| `services/api/internal/catalog/domain/food.go` | tipos FoodItem, NutritionFacts, HouseholdMeasure |
| `services/api/internal/catalog/repository/postgres.go` | CRUD alimentos, nutrição, medidas caseiras |
| `services/api/internal/catalog/usecase/catalog.go` | lógica de catálogo |
| `services/api/internal/catalog/handler.go` | endpoints /foods |

### Migrations (000016–000031)
| Arquivo | Conteúdo |
|---------|----------|
| `000016_create_professionals` | tabela professionals |
| `000017_create_professional_addresses` | endereços de atendimento |
| `000018_create_professional_service_modes` | modalidades (presencial/online) |
| `000019_create_availability_rules` | regras de disponibilidade |
| `000020_create_schedule_blocks` | bloqueios de agenda |
| `000021_create_appointments` | agendamentos |
| `000022_create_appointment_audit_events` | audit trail de agendamentos |
| `000023_create_patients` | tabela patients |
| `000024_create_patient_invites` | convites de paciente |
| `000025_create_patient_access_links` | links de acesso do paciente |
| `000026_create_patient_profiles` | perfil clínico do paciente |
| `000027_create_anamneses` | anamneses clínicas |
| `000028_create_progress_notes` | notas de evolução |
| `000029_create_clinical_attachments` | anexos clínicos |
| `000030_create_body_measurements` | medições de bioimpedância (20+ métricas) |
| `000031_create_food_catalog` | food_items, nutrition_facts, household_measures |

---

## Endpoints implementados

### Professionals (10 rotas)
- `POST/GET /professionals` — listar e criar
- `GET/PUT /professionals/{id}` — detalhe e atualização
- `POST/GET /professionals/{id}/addresses` — endereços
- `PUT /professionals/{id}/addresses/{addr_id}` — atualizar endereço
- `PUT/GET /professionals/{id}/service-modes` — modalidades

### Scheduling (12 rotas)
- `POST/GET /professionals/{prof_id}/availability` — regras de disponibilidade
- `PUT/DELETE /availability/{id}` — standalone
- `GET /professionals/{prof_id}/slots` — slots disponíveis
- `POST/GET /professionals/{prof_id}/blocks` — bloqueios
- `DELETE /blocks/{id}` — standalone
- `POST/GET /appointments` — agendamentos
- `GET /appointments/{id}` — detalhe
- `PATCH /appointments/{id}/status` — transição de status
- `PATCH /appointments/{id}/reschedule` — remarcação

### Patients (7 rotas)
- `POST/GET /patients` — listar e criar
- `GET/PUT /patients/{patient_id}` — detalhe e atualização
- `POST /patients/{patient_id}/invites` — gerar convite
- `POST/GET /patients/{patient_id}/profiles` — perfil clínico
- `POST /invites/activate` — ativação de acesso (auth-only)

### Clinical (7 rotas)
- `POST/GET /patients/{patient_id}/anamneses` — anamneses
- `GET/PUT /anamneses/{id}` — standalone
- `POST /anamneses/{id}/finalize` — finalização irreversível
- `POST/GET /patients/{patient_id}/notes` — notas de evolução
- `POST/GET /patients/{patient_id}/attachments` — anexos

### Bioimpedance (4 rotas)
- `POST /measurements` — criar medição
- `GET /measurements/{id}` — detalhe
- `POST /measurements/{id}/publish` — publicar
- `GET /patients/{patient_id}/measurements` — listar por paciente

### Food Catalog (6 rotas)
- `GET /foods` — listar alimentos
- `GET /foods/groups` — grupos alimentares
- `POST /foods` — criar alimento
- `GET/PUT/DELETE /foods/{id}` — CRUD individual

---

## Verificação

- [x] `go build ./...` compila sem erros
- [x] `go test ./...` — 10 pacotes ok, 0 falhas
- [x] Docker Compose: API sobe com todos os módulos
- [x] Migrations aplicam em banco limpo
- [x] Endpoints retornam 401 sem token
