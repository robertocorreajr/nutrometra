# Fase 2 — Operacao Clinica: Plano de Implementacao

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar os modulos de operacao clinica sobre a fundacao da Fase 1: profissionais/equipe/enderecos, disponibilidade/agenda, pacientes/convites, prontuario (anamnese/evolucao/anexos), bioimpedancia e catalogo alimentar — tudo tenant-scoped, auditado e protegido por RBAC.

**Architecture:** Novos modulos Go em `internal/{professional,scheduling,patient,clinical,bioimpedance,catalog}`, cada um com domain/repository/usecase/handler. Migrations 000016-000032 (16 schema + 1 seed). Todos os endpoints protegidos por AuthMiddleware + TenantMiddleware + RequirePermission. Operacoes mutaveis geram audit_logs. Entitlements checados onde aplicavel (limites de pacientes, profissionais, agendas por plano).

**Tech Stack:** Mesmo da Fase 1 — Go 1.26, chi/v5, pgx/v5, go-redis/v9, golang-migrate/v4, testify/v1, slog, Docker Compose, Zitadel v2.54.3.

**Pre-requisitos:** Fase 1 completa (tag v0.1.0-fase1), Docker stack funcional (`make setup`), migrations 000001-000015 aplicadas.

---

## Mapa de arquivos

### Migrations (000016-000031)
| Arquivo | Conteudo |
|---------|----------|
| `000016_create_professionals.{up,down}.sql` | tabela professionals |
| `000017_create_professional_addresses.{up,down}.sql` | tabela professional_addresses |
| `000018_create_professional_service_modes.{up,down}.sql` | tabela professional_service_modes |
| `000019_create_availability_rules.{up,down}.sql` | tabela availability_rules |
| `000020_create_schedule_blocks.{up,down}.sql` | tabela schedule_blocks |
| `000021_create_appointments.{up,down}.sql` | tabela appointments |
| `000022_create_appointment_audit_events.{up,down}.sql` | tabela appointment_audit_events |
| `000023_create_patients.{up,down}.sql` | tabela patients |
| `000024_create_patient_invites.{up,down}.sql` | tabela patient_invites |
| `000025_create_patient_access_links.{up,down}.sql` | tabela patient_access_links |
| `000026_create_patient_profiles.{up,down}.sql` | tabela patient_profiles + patient_consents |
| `000027_create_anamneses.{up,down}.sql` | tabela anamneses |
| `000028_create_progress_notes.{up,down}.sql` | tabela progress_notes |
| `000029_create_clinical_attachments.{up,down}.sql` | tabela clinical_attachments |
| `000030_create_body_measurements.{up,down}.sql` | tabela body_measurements + body_measurement_publications |
| `000031_create_food_catalog.{up,down}.sql` | tabelas food_items, food_nutrition_facts, food_household_measures |

### Go API — professional
| Arquivo | Responsabilidade |
|---------|-----------------|
| `internal/professional/domain/professional.go` | entidades Professional, Address, ServiceMode |
| `internal/professional/repository/postgres.go` | CRUD professionals, addresses, service_modes |
| `internal/professional/usecase/professional.go` | logica de negocio, validacoes CRN |
| `internal/professional/handler.go` | endpoints /professionals, /professionals/{id}/addresses |

### Go API — scheduling
| Arquivo | Responsabilidade |
|---------|-----------------|
| `internal/scheduling/domain/scheduling.go` | entidades AvailabilityRule, ScheduleBlock, Appointment |
| `internal/scheduling/repository/postgres.go` | CRUD availability, blocks, appointments |
| `internal/scheduling/usecase/scheduling.go` | validacao de conflitos, slots disponiveis, buffers |
| `internal/scheduling/handler.go` | endpoints /availability, /blocks, /appointments |

### Go API — patient
| Arquivo | Responsabilidade |
|---------|-----------------|
| `internal/patient/domain/patient.go` | entidades Patient, PatientInvite, PatientAccessLink, PatientProfile |
| `internal/patient/repository/postgres.go` | CRUD patients, invites, access_links, profiles |
| `internal/patient/usecase/patient.go` | cadastro, convites, ativacao portal |
| `internal/patient/handler.go` | endpoints /patients, /patients/{id}/invites |

### Go API — clinical
| Arquivo | Responsabilidade |
|---------|-----------------|
| `internal/clinical/domain/clinical.go` | entidades Anamnesis, ProgressNote, ClinicalAttachment |
| `internal/clinical/repository/postgres.go` | CRUD anamneses, progress_notes, attachments |
| `internal/clinical/usecase/clinical.go` | rascunho/finalizacao, publicacao ao paciente |
| `internal/clinical/handler.go` | endpoints /patients/{id}/anamneses, /patients/{id}/notes |

### Go API — bioimpedance
| Arquivo | Responsabilidade |
|---------|-----------------|
| `internal/bioimpedance/domain/measurement.go` | entidades BodyMeasurement, Publication |
| `internal/bioimpedance/repository/postgres.go` | CRUD body_measurements, publications |
| `internal/bioimpedance/usecase/bioimpedance.go` | validacao de ranges, historico, publicacao |
| `internal/bioimpedance/handler.go` | endpoints /patients/{id}/measurements |

### Go API — catalog
| Arquivo | Responsabilidade |
|---------|-----------------|
| `internal/catalog/domain/food.go` | entidades FoodItem, NutritionFacts, HouseholdMeasure |
| `internal/catalog/repository/postgres.go` | busca, CRUD tenant-scoped, listagem global |
| `internal/catalog/usecase/catalog.go` | busca com filtro, criacao tenant-scoped |
| `internal/catalog/handler.go` | endpoints /foods, /foods/{id}/measures |

---

## Task 1: Branch e migrations — professionals, addresses, service_modes

**Files:**
- Create: `services/api/migrations/000016_create_professionals.up.sql`
- Create: `services/api/migrations/000016_create_professionals.down.sql`
- Create: `services/api/migrations/000017_create_professional_addresses.up.sql`
- Create: `services/api/migrations/000017_create_professional_addresses.down.sql`
- Create: `services/api/migrations/000018_create_professional_service_modes.up.sql`
- Create: `services/api/migrations/000018_create_professional_service_modes.down.sql`

- [ ] **Step 1: Criar branch**

```bash
git checkout -b feat/fase-2-operacao-clinica
```

- [ ] **Step 2: Migration 000016 — professionals**

`services/api/migrations/000016_create_professionals.up.sql`:
```sql
CREATE TABLE professionals (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    user_id     UUID REFERENCES users(id),
    full_name   TEXT NOT NULL,
    registration_type TEXT NOT NULL DEFAULT 'CRN'
        CHECK (registration_type IN ('CRN', 'CRO', 'CRM', 'other')),
    registration_number TEXT,
    registration_state  TEXT,
    document_number TEXT,
    phone       TEXT,
    email       TEXT,
    specialty_tags_json JSONB DEFAULT '[]'::jsonb,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT professionals_crn_state_required
        CHECK (registration_type != 'CRN' OR registration_state IS NOT NULL)
);

CREATE INDEX idx_professionals_tenant ON professionals(tenant_id);
CREATE INDEX idx_professionals_user   ON professionals(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_professionals_tenant_active ON professionals(tenant_id, is_active) WHERE is_active = true;
```

`services/api/migrations/000016_create_professionals.down.sql`:
```sql
DROP TABLE IF EXISTS professionals;
```

- [ ] **Step 3: Migration 000017 — professional_addresses**

`services/api/migrations/000017_create_professional_addresses.up.sql`:
```sql
CREATE TABLE professional_addresses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    label       TEXT NOT NULL,
    street      TEXT NOT NULL,
    number      TEXT NOT NULL,
    complement  TEXT,
    district    TEXT NOT NULL,
    city        TEXT NOT NULL,
    state       TEXT NOT NULL,
    zip_code    TEXT NOT NULL,
    country     TEXT NOT NULL DEFAULT 'BR',
    latitude    DOUBLE PRECISION,
    longitude   DOUBLE PRECISION,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_prof_addresses_tenant ON professional_addresses(tenant_id);
CREATE INDEX idx_prof_addresses_prof   ON professional_addresses(professional_id);
```

`services/api/migrations/000017_create_professional_addresses.down.sql`:
```sql
DROP TABLE IF EXISTS professional_addresses;
```

- [ ] **Step 4: Migration 000018 — professional_service_modes**

`services/api/migrations/000018_create_professional_service_modes.up.sql`:
```sql
CREATE TABLE professional_service_modes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    mode        TEXT NOT NULL CHECK (mode IN ('onsite', 'online')),
    is_enabled  BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (professional_id, mode)
);

CREATE INDEX idx_prof_service_modes_prof ON professional_service_modes(professional_id);
```

`services/api/migrations/000018_create_professional_service_modes.down.sql`:
```sql
DROP TABLE IF EXISTS professional_service_modes;
```

- [ ] **Step 5: Rodar migrations**

```bash
make migrate-up
```

- [ ] **Step 6: Commit**

```bash
git add services/api/migrations/00001{6,7,8}_*
git commit -m "feat(phase2): add migrations for professionals, addresses, service_modes"
```

---

## Task 2: Professional domain, repository, usecase e handler

**Files:**
- Create: `services/api/internal/professional/domain/professional.go`
- Create: `services/api/internal/professional/repository/postgres.go`
- Create: `services/api/internal/professional/usecase/professional.go`
- Create: `services/api/internal/professional/handler.go`

- [ ] **Step 1: Domain — entidades e erros**

`services/api/internal/professional/domain/professional.go`:
- Professional struct (id, tenant_id, user_id, full_name, registration_type/number/state, document_number, phone, email, specialty_tags, is_active, created_at, updated_at)
- Address struct (id, tenant_id, professional_id, label, street, number, complement, district, city, state, zip_code, country, lat, lng, is_active)
- ServiceMode struct (id, tenant_id, professional_id, mode, is_enabled)
- RegistrationType constants: CRN, CRO, CRM, Other
- Sentinel errors: ErrProfessionalNotFound, ErrAddressNotFound, ErrDuplicateServiceMode
- Validate() method on Professional — registration_state required when CRN

- [ ] **Step 2: Repository — postgres CRUD**

`services/api/internal/professional/repository/postgres.go`:
- Create(ctx, Professional) (Professional, error) — INSERT + RETURNING
- GetByID(ctx, tenantID, id) (Professional, error)
- ListByTenant(ctx, tenantID, onlyActive bool) ([]Professional, error)
- Update(ctx, Professional) error
- CreateAddress(ctx, Address) (Address, error)
- ListAddresses(ctx, tenantID, professionalID) ([]Address, error)
- UpdateAddress(ctx, Address) error
- SetServiceMode(ctx, ServiceMode) error — UPSERT
- ListServiceModes(ctx, tenantID, professionalID) ([]ServiceMode, error)
- Todos os queries filtram por tenant_id

- [ ] **Step 3: Usecase — logica de negocio**

`services/api/internal/professional/usecase/professional.go`:
- ProfessionalUsecase struct com repo, auditSvc, entitlementSvc
- CreateProfessional — valida dominio, checa entitlement `max_professionals`, cria, audita
- GetByID — delega ao repo
- ListByTenant — delega ao repo
- UpdateProfessional — valida, atualiza, audita
- CreateAddress — checa entitlement `max_addresses`, cria, audita
- ListAddresses — delega
- UpdateAddress — atualiza, audita
- SetServiceMode — upsert, audita

- [ ] **Step 4: Handler — endpoints HTTP**

`services/api/internal/professional/handler.go`:
- POST   /professionals — RequirePermission("tenant:manage") — CreateProfessional
- GET    /professionals — RequirePermission("patients:read") — ListByTenant
- GET    /professionals/{id} — RequirePermission("patients:read") — GetByID
- PUT    /professionals/{id} — RequirePermission("tenant:manage") — UpdateProfessional
- POST   /professionals/{id}/addresses — RequirePermission("tenant:manage") — CreateAddress
- GET    /professionals/{id}/addresses — RequirePermission("schedule:read") — ListAddresses
- PUT    /professionals/{id}/addresses/{addr_id} — RequirePermission("tenant:manage") — UpdateAddress
- PUT    /professionals/{id}/service-modes — RequirePermission("tenant:manage") — SetServiceMode
- GET    /professionals/{id}/service-modes — RequirePermission("schedule:read") — ListServiceModes

Request/response DTOs separados do dominio. Validacao de input no handler. UUID parsing com erro 400.

- [ ] **Step 5: Testes unitarios**

- domain/professional_test.go — Validate() com CRN sem estado falha, com estado passa
- usecase/professional_test.go — mock repo, testar entitlement check, audit call

- [ ] **Step 6: Wire-up em main.go**

Registrar professional handler no router, dentro do grupo tenant-scoped autenticado.

- [ ] **Step 7: Commit**

```bash
git add services/api/internal/professional/ services/api/cmd/api/main.go
git commit -m "feat(phase2): professional module — domain, repository, usecase, handler"
```

---

## Task 3: Migrations — availability_rules, schedule_blocks

**Files:**
- Create: `services/api/migrations/000019_create_availability_rules.up.sql`
- Create: `services/api/migrations/000019_create_availability_rules.down.sql`
- Create: `services/api/migrations/000020_create_schedule_blocks.up.sql`
- Create: `services/api/migrations/000020_create_schedule_blocks.down.sql`

- [ ] **Step 1: Migration 000019 — availability_rules**

```sql
CREATE TABLE availability_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    address_id      UUID REFERENCES professional_addresses(id),
    service_mode    TEXT NOT NULL CHECK (service_mode IN ('onsite', 'online')),
    weekday         SMALLINT NOT NULL CHECK (weekday BETWEEN 0 AND 6),
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    slot_minutes    INT NOT NULL DEFAULT 50,
    buffer_before_minutes INT NOT NULL DEFAULT 0,
    buffer_after_minutes  INT NOT NULL DEFAULT 10,
    recurrence_end_date DATE,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT avail_end_after_start CHECK (end_time > start_time),
    CONSTRAINT avail_onsite_needs_address
        CHECK (service_mode != 'onsite' OR address_id IS NOT NULL)
);

CREATE INDEX idx_avail_rules_prof ON availability_rules(professional_id, weekday);
CREATE INDEX idx_avail_rules_tenant ON availability_rules(tenant_id);
```

- [ ] **Step 2: Migration 000020 — schedule_blocks**

```sql
CREATE TABLE schedule_blocks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    address_id      UUID REFERENCES professional_addresses(id),
    service_mode    TEXT CHECK (service_mode IN ('onsite', 'online')),
    starts_at       TIMESTAMPTZ NOT NULL,
    ends_at         TIMESTAMPTZ NOT NULL,
    reason          TEXT NOT NULL,
    created_by_user_id UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT block_end_after_start CHECK (ends_at > starts_at)
);

CREATE INDEX idx_schedule_blocks_prof ON schedule_blocks(professional_id, starts_at, ends_at);
CREATE INDEX idx_schedule_blocks_tenant ON schedule_blocks(tenant_id);
```

- [ ] **Step 3: Rodar migrations e commit**

```bash
make migrate-up
git add services/api/migrations/00001{9,20}_*
git commit -m "feat(phase2): add migrations for availability_rules and schedule_blocks"
```

---

## Task 4: Migration — appointments e appointment_audit_events

**Files:**
- Create: `services/api/migrations/000021_create_appointments.up.sql`
- Create: `services/api/migrations/000021_create_appointments.down.sql`
- Create: `services/api/migrations/000022_create_appointment_audit_events.up.sql`
- Create: `services/api/migrations/000022_create_appointment_audit_events.down.sql`

- [ ] **Step 1: Migration 000021 — appointments**

```sql
CREATE TABLE appointments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    patient_id      UUID NOT NULL, -- FK added after patients table exists
    address_id      UUID REFERENCES professional_addresses(id),
    service_mode    TEXT NOT NULL CHECK (service_mode IN ('onsite', 'online')),
    source          TEXT NOT NULL DEFAULT 'manual'
        CHECK (source IN ('manual', 'patient_request', 'import', 'sync')),
    status          TEXT NOT NULL DEFAULT 'scheduled'
        CHECK (status IN ('scheduled', 'confirmed', 'completed', 'cancelled', 'no_show')),
    scheduled_start_at TIMESTAMPTZ NOT NULL,
    scheduled_end_at   TIMESTAMPTZ NOT NULL,
    cancellation_reason TEXT,
    notes           TEXT,
    external_calendar_event_id TEXT,
    created_by_user_id  UUID NOT NULL REFERENCES users(id),
    updated_by_user_id  UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT appt_end_after_start CHECK (scheduled_end_at > scheduled_start_at),
    CONSTRAINT appt_onsite_needs_address
        CHECK (service_mode != 'onsite' OR address_id IS NOT NULL)
);

CREATE INDEX idx_appointments_prof_time ON appointments(professional_id, scheduled_start_at, scheduled_end_at);
CREATE INDEX idx_appointments_patient ON appointments(patient_id);
CREATE INDEX idx_appointments_tenant ON appointments(tenant_id);
CREATE INDEX idx_appointments_status ON appointments(tenant_id, status);
```

- [ ] **Step 2: Migration 000022 — appointment_audit_events**

```sql
CREATE TABLE appointment_audit_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    appointment_id  UUID NOT NULL REFERENCES appointments(id),
    event_type      TEXT NOT NULL,
    actor_user_id   UUID REFERENCES users(id),
    payload_json    JSONB DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_appt_audit_appointment ON appointment_audit_events(appointment_id);
CREATE INDEX idx_appt_audit_tenant ON appointment_audit_events(tenant_id);
```

- [ ] **Step 3: Rodar migrations e commit**

```bash
make migrate-up
git add services/api/migrations/00002{1,2}_*
git commit -m "feat(phase2): add migrations for appointments and appointment_audit_events"
```

---

## Task 5: Scheduling module — domain, repository, usecase, handler

**Files:**
- Create: `services/api/internal/scheduling/domain/scheduling.go`
- Create: `services/api/internal/scheduling/repository/postgres.go`
- Create: `services/api/internal/scheduling/usecase/scheduling.go`
- Create: `services/api/internal/scheduling/handler.go`

- [ ] **Step 1: Domain — entidades**

`services/api/internal/scheduling/domain/scheduling.go`:
- AvailabilityRule struct — todos os campos da tabela
- ScheduleBlock struct
- Appointment struct
- AppointmentAuditEvent struct
- AppointmentStatus constants: Scheduled, Confirmed, Completed, Cancelled, NoShow
- ServiceMode constants: Onsite, Online
- Source constants: Manual, PatientRequest, Import, Sync
- TimeSlot struct { Start, End time.Time; Available bool } — para calculo de slots
- Sentinel errors: ErrAppointmentNotFound, ErrConflict, ErrBlockedTime, ErrOutsideAvailability
- Validate() nos structs relevantes

- [ ] **Step 2: Repository**

`services/api/internal/scheduling/repository/postgres.go`:
- AvailabilityRule CRUD (Create, List by professional, Update, Delete)
- ScheduleBlock CRUD (Create, List by professional+range, Delete)
- Appointment CRUD (Create, GetByID, List by professional+range, List by patient, Update status)
- FindConflictingAppointments(ctx, tenantID, professionalID, start, end, excludeID) — para validacao
- CreateAuditEvent(ctx, event) — registra mudanca no appointment
- Todos filtram por tenant_id

- [ ] **Step 3: Usecase — logica de agendamento**

`services/api/internal/scheduling/usecase/scheduling.go`:
- SchedulingUsecase com repo, auditSvc, entitlementSvc
- **Availability:**
  - CreateAvailabilityRule — valida, cria, audita
  - ListAvailabilityRules — por profissional
  - UpdateAvailabilityRule — valida, atualiza
  - DeleteAvailabilityRule — soft delete (is_active = false)
  - GetAvailableSlots(professionalID, date, serviceMode) — calcula slots livres considerando rules, blocks, appointments existentes
- **Blocks:**
  - CreateBlock — valida, cria, audita
  - ListBlocks — por profissional e range
  - DeleteBlock — remove, audita
- **Appointments:**
  - CreateAppointment — valida disponibilidade, checa conflitos, checa entitlement, cria, registra audit_event, audita
  - GetAppointment — por ID
  - ListAppointments — por profissional ou paciente, com filtro de range e status
  - UpdateAppointmentStatus — transicoes validas (scheduled->confirmed->completed, scheduled->cancelled, confirmed->cancelled, *->no_show), registra audit_event
  - RescheduleAppointment — valida novo horario, atualiza, registra audit_event
  - CancelAppointment — requer motivo, registra audit_event

- [ ] **Step 4: Handler — endpoints HTTP**

`services/api/internal/scheduling/handler.go`:
- **Availability:**
  - POST   /professionals/{prof_id}/availability — RequirePermission("schedule:manage")
  - GET    /professionals/{prof_id}/availability — RequirePermission("schedule:read")
  - PUT    /availability/{id} — RequirePermission("schedule:manage")
  - DELETE /availability/{id} — RequirePermission("schedule:manage")
  - GET    /professionals/{prof_id}/slots?date=YYYY-MM-DD&mode=onsite|online — RequirePermission("schedule:read")
- **Blocks:**
  - POST   /professionals/{prof_id}/blocks — RequirePermission("schedule:manage")
  - GET    /professionals/{prof_id}/blocks?from=&to= — RequirePermission("schedule:read")
  - DELETE /blocks/{id} — RequirePermission("schedule:manage")
- **Appointments:**
  - POST   /appointments — RequirePermission("schedule:manage")
  - GET    /appointments?professional_id=&patient_id=&from=&to=&status= — RequirePermission("appointments:read")
  - GET    /appointments/{id} — RequirePermission("appointments:read")
  - PATCH  /appointments/{id}/status — RequirePermission("schedule:manage")
  - PATCH  /appointments/{id}/reschedule — RequirePermission("schedule:manage")

- [ ] **Step 5: Testes unitarios**

- domain/scheduling_test.go — validacao de entidades
- usecase/scheduling_test.go — mock repo, testar deteccao de conflitos, transicoes de status invalidas, calculo de slots

- [ ] **Step 6: Wire-up em main.go**

Registrar scheduling handler no router.

- [ ] **Step 7: Commit**

```bash
git add services/api/internal/scheduling/ services/api/cmd/api/main.go
git commit -m "feat(phase2): scheduling module — availability, blocks, appointments"
```

---

## Task 6: Migrations — patients, invites, access_links, profiles, consents

**Files:**
- Create: `services/api/migrations/000023_create_patients.{up,down}.sql`
- Create: `services/api/migrations/000024_create_patient_invites.{up,down}.sql`
- Create: `services/api/migrations/000025_create_patient_access_links.{up,down}.sql`
- Create: `services/api/migrations/000026_create_patient_profiles.{up,down}.sql`

- [ ] **Step 1: Migration 000023 — patients**

```sql
CREATE TABLE patients (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants(id),
    primary_professional_id UUID NOT NULL REFERENCES professionals(id),
    external_code           TEXT,
    full_name               TEXT NOT NULL,
    preferred_name          TEXT,
    birth_date              DATE,
    sex                     TEXT CHECK (sex IN ('M', 'F', 'other')),
    phone                   TEXT,
    email                   TEXT,
    contact_json            JSONB DEFAULT '{}'::jsonb,
    status                  TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'inactive', 'archived')),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_patients_tenant ON patients(tenant_id);
CREATE INDEX idx_patients_professional ON patients(primary_professional_id);
CREATE INDEX idx_patients_tenant_status ON patients(tenant_id, status);
CREATE INDEX idx_patients_name ON patients(tenant_id, full_name);

-- Add FK from appointments to patients (deferred from migration 000021)
ALTER TABLE appointments ADD CONSTRAINT fk_appointments_patient
    FOREIGN KEY (patient_id) REFERENCES patients(id);
```

- [ ] **Step 2: Migration 000024 — patient_invites**

```sql
CREATE TABLE patient_invites (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id),
    patient_id          UUID REFERENCES patients(id),
    professional_id     UUID NOT NULL REFERENCES professionals(id),
    invite_code_hash    TEXT NOT NULL,
    expires_at          TIMESTAMPTZ NOT NULL,
    max_uses            INT NOT NULL DEFAULT 1,
    used_count          INT NOT NULL DEFAULT 0,
    created_by_user_id  UUID NOT NULL REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_patient_invites_code ON patient_invites(invite_code_hash);
CREATE INDEX idx_patient_invites_tenant ON patient_invites(tenant_id);
```

- [ ] **Step 3: Migration 000025 — patient_access_links**

```sql
CREATE TABLE patient_access_links (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    patient_id  UUID NOT NULL REFERENCES patients(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    activated_at TIMESTAMPTZ,
    status      TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'active', 'revoked')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (tenant_id, patient_id, user_id)
);

CREATE INDEX idx_patient_access_tenant ON patient_access_links(tenant_id);
CREATE INDEX idx_patient_access_user ON patient_access_links(user_id);
```

- [ ] **Step 4: Migration 000026 — patient_profiles e patient_consents**

```sql
CREATE TABLE patient_profiles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    patient_id  UUID NOT NULL REFERENCES patients(id),
    profile_type TEXT NOT NULL,
    priority    INT NOT NULL DEFAULT 0,
    metadata_json JSONB DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_patient_profiles_patient ON patient_profiles(patient_id);

CREATE TABLE patient_consents (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    patient_id  UUID NOT NULL REFERENCES patients(id),
    consent_type TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'accepted', 'revoked')),
    accepted_at TIMESTAMPTZ,
    revoked_at  TIMESTAMPTZ,
    metadata_json JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX idx_patient_consents_patient ON patient_consents(patient_id);
```

- [ ] **Step 5: Rodar migrations e commit**

```bash
make migrate-up
git add services/api/migrations/00002{3,4,5,6}_*
git commit -m "feat(phase2): add migrations for patients, invites, access_links, profiles, consents"
```

---

## Task 7: Patient module — domain, repository, usecase, handler

**Files:**
- Create: `services/api/internal/patient/domain/patient.go`
- Create: `services/api/internal/patient/repository/postgres.go`
- Create: `services/api/internal/patient/usecase/patient.go`
- Create: `services/api/internal/patient/handler.go`

- [ ] **Step 1: Domain**

`services/api/internal/patient/domain/patient.go`:
- Patient struct com todos os campos
- PatientInvite struct — invite_code_hash, expires_at, max_uses, used_count
- PatientAccessLink struct — user_id, status
- PatientProfile struct — profile_type (adulto_geral, infantil, tea, gestante, lactante, idoso, atleta, fisiculturista, emagrecimento, diabetes, gastro, vegetariano, comportamento_alimentar)
- PatientConsent struct — consent_type, status
- PatientStatus: Active, Inactive, Archived
- Sentinel errors: ErrPatientNotFound, ErrInviteExpired, ErrInviteMaxUses, ErrInviteNotFound

- [ ] **Step 2: Repository**

- Create(ctx, Patient) (Patient, error)
- GetByID(ctx, tenantID, id)
- ListByTenant(ctx, tenantID, filters) ([]Patient, int, error) — com paginacao
- ListByProfessional(ctx, tenantID, professionalID) ([]Patient, error)
- Update(ctx, Patient) error
- CreateInvite(ctx, PatientInvite) (PatientInvite, error)
- GetInviteByCodeHash(ctx, hash) (PatientInvite, error)
- IncrementInviteUsage(ctx, inviteID) error
- CreateAccessLink(ctx, PatientAccessLink) error
- GetAccessLinkByUser(ctx, tenantID, userID) (PatientAccessLink, error)
- CreateProfile(ctx, PatientProfile) error
- ListProfiles(ctx, tenantID, patientID) ([]PatientProfile, error)
- CreateConsent(ctx, PatientConsent) error

- [ ] **Step 3: Usecase**

- CreatePatient — valida, checa entitlement `max_patients`, cria, audita
- GetByID, ListByTenant, ListByProfessional — delegam
- UpdatePatient — valida, atualiza, audita
- GenerateInvite — gera codigo aleatorio, hasheia (SHA-256), cria invite com expiracao, audita
- ActivatePortalAccess(inviteCode, userID) — valida invite (expiracao, max_uses), cria access_link, incrementa used_count, audita
- AddProfile — cria perfil nutricional
- RecordConsent — registra consentimento

- [ ] **Step 4: Handler**

- POST   /patients — RequirePermission("patients:write")
- GET    /patients — RequirePermission("patients:read") — com query params: ?status=&professional_id=&q=&page=&per_page=
- GET    /patients/{id} — RequirePermission("patients:read")
- PUT    /patients/{id} — RequirePermission("patients:write")
- POST   /patients/{id}/invites — RequirePermission("patients:write") — gera invite code
- POST   /invites/activate — publico autenticado (sem RequirePermission especifico, apenas auth) — ativa portal
- POST   /patients/{id}/profiles — RequirePermission("patients:write")
- GET    /patients/{id}/profiles — RequirePermission("patients:read")

- [ ] **Step 5: Testes**

- domain/patient_test.go — validacao de entidades
- usecase/patient_test.go — mock, testar invite flow (gerar, ativar, expirado, max_uses)

- [ ] **Step 6: Wire-up em main.go**

- [ ] **Step 7: Commit**

```bash
git add services/api/internal/patient/ services/api/cmd/api/main.go
git commit -m "feat(phase2): patient module — domain, repository, usecase, handler with invite flow"
```

---

## Task 8: Migrations — anamneses, progress_notes, clinical_attachments

**Files:**
- Create: `services/api/migrations/000027_create_anamneses.{up,down}.sql`
- Create: `services/api/migrations/000028_create_progress_notes.{up,down}.sql`
- Create: `services/api/migrations/000029_create_clinical_attachments.{up,down}.sql`

- [ ] **Step 1: Migration 000027 — anamneses**

```sql
CREATE TABLE anamneses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    appointment_id  UUID REFERENCES appointments(id),
    status          TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'finalized')),
    form_version    INT NOT NULL DEFAULT 1,
    payload_json    JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_anamneses_patient ON anamneses(patient_id);
CREATE INDEX idx_anamneses_tenant ON anamneses(tenant_id);
```

- [ ] **Step 2: Migration 000028 — progress_notes**

```sql
CREATE TABLE progress_notes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    appointment_id  UUID REFERENCES appointments(id),
    note_type       TEXT NOT NULL DEFAULT 'evolution'
        CHECK (note_type IN ('evolution', 'observation', 'followup')),
    content_json    JSONB NOT NULL DEFAULT '{}'::jsonb,
    visibility_to_patient TEXT NOT NULL DEFAULT 'hidden'
        CHECK (visibility_to_patient IN ('hidden', 'published')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_progress_notes_patient ON progress_notes(patient_id);
CREATE INDEX idx_progress_notes_tenant ON progress_notes(tenant_id);
```

- [ ] **Step 3: Migration 000029 — clinical_attachments**

```sql
CREATE TABLE clinical_attachments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    appointment_id  UUID REFERENCES appointments(id),
    file_key        TEXT NOT NULL,
    file_name       TEXT NOT NULL,
    mime_type       TEXT NOT NULL,
    size_bytes      BIGINT NOT NULL,
    category        TEXT NOT NULL DEFAULT 'general',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_clinical_attach_patient ON clinical_attachments(patient_id);
CREATE INDEX idx_clinical_attach_tenant ON clinical_attachments(tenant_id);
```

- [ ] **Step 4: Rodar migrations e commit**

```bash
make migrate-up
git add services/api/migrations/00002{7,8,9}_*
git commit -m "feat(phase2): add migrations for anamneses, progress_notes, clinical_attachments"
```

---

## Task 9: Clinical module — domain, repository, usecase, handler

**Files:**
- Create: `services/api/internal/clinical/domain/clinical.go`
- Create: `services/api/internal/clinical/repository/postgres.go`
- Create: `services/api/internal/clinical/usecase/clinical.go`
- Create: `services/api/internal/clinical/handler.go`

- [ ] **Step 1: Domain**

- Anamnesis struct — status (draft/finalized), form_version, payload_json (JSONB flexivel para formularios)
- ProgressNote struct — note_type (evolution/observation/followup), content_json, visibility_to_patient
- ClinicalAttachment struct — file_key, file_name, mime_type, size_bytes, category
- Sentinel errors: ErrAnamnesisNotFound, ErrNoteNotFound, ErrAttachmentNotFound, ErrAlreadyFinalized

- [ ] **Step 2: Repository**

- Anamnesis: Create, GetByID, ListByPatient, Update, Finalize (set status=finalized, bloqueia edicao)
- ProgressNote: Create, GetByID, ListByPatient, Update, SetVisibility
- ClinicalAttachment: Create, GetByID, ListByPatient, Delete
- Todos filtram por tenant_id

- [ ] **Step 3: Usecase**

- CreateAnamnesis — valida, cria como draft, audita
- UpdateAnamnesis — verifica draft, atualiza payload, audita
- FinalizeAnamnesis — muda status para finalized (irreversivel), audita
- GetAnamnesis, ListByPatient
- CreateProgressNote — cria, audita
- UpdateProgressNote — atualiza, audita
- PublishNoteToPatient — muda visibility, audita
- CreateAttachment — registra metadata (upload real via signed URL futuro), audita
- ListAttachments, DeleteAttachment — audita

- [ ] **Step 4: Handler**

- POST   /patients/{patient_id}/anamneses — RequirePermission("clinical:write")
- GET    /patients/{patient_id}/anamneses — RequirePermission("clinical:read")
- GET    /anamneses/{id} — RequirePermission("clinical:read")
- PUT    /anamneses/{id} — RequirePermission("clinical:write")
- POST   /anamneses/{id}/finalize — RequirePermission("clinical:write")
- POST   /patients/{patient_id}/notes — RequirePermission("clinical:write")
- GET    /patients/{patient_id}/notes — RequirePermission("clinical:read")
- PUT    /notes/{id} — RequirePermission("clinical:write")
- PATCH  /notes/{id}/visibility — RequirePermission("clinical:write")
- POST   /patients/{patient_id}/attachments — RequirePermission("clinical:write")
- GET    /patients/{patient_id}/attachments — RequirePermission("clinical:read")
- DELETE /attachments/{id} — RequirePermission("clinical:write")

- [ ] **Step 5: Testes**

- usecase/clinical_test.go — testar finalizacao irreversivel, publicacao de nota

- [ ] **Step 6: Wire-up em main.go**

- [ ] **Step 7: Commit**

```bash
git add services/api/internal/clinical/ services/api/cmd/api/main.go
git commit -m "feat(phase2): clinical module — anamnesis, progress notes, attachments"
```

---

## Task 10: Migrations — body_measurements e publications

**Files:**
- Create: `services/api/migrations/000030_create_body_measurements.{up,down}.sql`

- [ ] **Step 1: Migration 000030 — body_measurements e body_measurement_publications**

```sql
CREATE TABLE body_measurements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    appointment_id  UUID REFERENCES appointments(id),
    measured_at     TIMESTAMPTZ NOT NULL,
    source          TEXT NOT NULL DEFAULT 'manual'
        CHECK (source IN ('manual', 'imported')),
    weight_kg               DOUBLE PRECISION,
    height_cm               DOUBLE PRECISION,
    bmi                     DOUBLE PRECISION,
    body_fat_percent        DOUBLE PRECISION,
    skeletal_muscle_percent DOUBLE PRECISION,
    muscle_mass_kg          DOUBLE PRECISION,
    fat_mass_kg             DOUBLE PRECISION,
    visceral_fat_index      DOUBLE PRECISION,
    body_water_percent      DOUBLE PRECISION,
    basal_metabolic_rate_kcal DOUBLE PRECISION,
    metabolic_age           INT,
    waist_cm                DOUBLE PRECISION,
    hip_cm                  DOUBLE PRECISION,
    notes                   TEXT,
    raw_payload_json        JSONB,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_body_meas_patient ON body_measurements(patient_id, measured_at DESC);
CREATE INDEX idx_body_meas_tenant ON body_measurements(tenant_id);

CREATE TABLE body_measurement_publications (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id),
    body_measurement_id UUID NOT NULL REFERENCES body_measurements(id),
    patient_id          UUID NOT NULL REFERENCES patients(id),
    published_by_user_id UUID NOT NULL REFERENCES users(id),
    published_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    visibility_mode     TEXT NOT NULL DEFAULT 'summary'
        CHECK (visibility_mode IN ('summary', 'full'))
);

CREATE INDEX idx_body_meas_pub_measurement ON body_measurement_publications(body_measurement_id);
```

- [ ] **Step 2: Rodar migration e commit**

```bash
make migrate-up
git add services/api/migrations/000030_*
git commit -m "feat(phase2): add migration for body_measurements and publications"
```

---

## Task 11: Bioimpedance module — domain, repository, usecase, handler

**Files:**
- Create: `services/api/internal/bioimpedance/domain/measurement.go`
- Create: `services/api/internal/bioimpedance/repository/postgres.go`
- Create: `services/api/internal/bioimpedance/usecase/bioimpedance.go`
- Create: `services/api/internal/bioimpedance/handler.go`

- [ ] **Step 1: Domain**

- BodyMeasurement struct — todos os campos da tabela
- MeasurementPublication struct
- VisibilityMode: Summary, Full
- Validate() — peso > 0 quando informado, altura > 0 quando informada, IMC calculavel se peso+altura presentes
- Sentinel errors: ErrMeasurementNotFound, ErrAlreadyPublished

- [ ] **Step 2: Repository**

- Create(ctx, BodyMeasurement) (BodyMeasurement, error)
- GetByID(ctx, tenantID, id)
- ListByPatient(ctx, tenantID, patientID) — ordenado por measured_at DESC
- GetHistory(ctx, tenantID, patientID, limit) — para comparacao temporal
- Publish(ctx, MeasurementPublication) error
- GetPublishedForPatient(ctx, tenantID, patientID) — medicoes publicadas

- [ ] **Step 3: Usecase**

- RecordMeasurement — valida ranges basicos (peso 1-500kg, altura 30-300cm), calcula BMI se possivel, cria, audita
- GetMeasurement, ListByPatient
- GetEvolutionHistory — retorna serie temporal para comparacao
- PublishToPatient — cria publication, audita
- GetPublishedMeasurements — para portal paciente

- [ ] **Step 4: Handler**

- POST   /patients/{patient_id}/measurements — RequirePermission("clinical:write")
- GET    /patients/{patient_id}/measurements — RequirePermission("clinical:read")
- GET    /measurements/{id} — RequirePermission("clinical:read")
- GET    /patients/{patient_id}/measurements/history?limit= — RequirePermission("clinical:read")
- POST   /measurements/{id}/publish — RequirePermission("clinical:write")

- [ ] **Step 5: Testes**

- domain/measurement_test.go — validacao de ranges, calculo BMI
- usecase/bioimpedance_test.go — mock repo, testar validacao e publicacao

- [ ] **Step 6: Wire-up em main.go**

- [ ] **Step 7: Commit**

```bash
git add services/api/internal/bioimpedance/ services/api/cmd/api/main.go
git commit -m "feat(phase2): bioimpedance module — measurements, history, publication"
```

---

## Task 12: Migration — food catalog

**Files:**
- Create: `services/api/migrations/000031_create_food_catalog.{up,down}.sql`

- [ ] **Step 1: Migration 000031 — food_items, food_nutrition_facts, food_household_measures**

```sql
CREATE TABLE food_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID REFERENCES tenants(id),
    source_scope TEXT NOT NULL DEFAULT 'global'
        CHECK (source_scope IN ('global', 'tenant')),
    name        TEXT NOT NULL,
    common_name TEXT,
    brand       TEXT,
    food_group  TEXT,
    enabled     BOOLEAN NOT NULL DEFAULT true,
    metadata_json JSONB DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT food_tenant_scope_check
        CHECK ((source_scope = 'global' AND tenant_id IS NULL) OR
               (source_scope = 'tenant' AND tenant_id IS NOT NULL))
);

CREATE INDEX idx_food_items_tenant ON food_items(tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX idx_food_items_name ON food_items(name);
CREATE INDEX idx_food_items_group ON food_items(food_group) WHERE food_group IS NOT NULL;
CREATE INDEX idx_food_items_global_enabled ON food_items(source_scope, enabled) WHERE source_scope = 'global';

CREATE TABLE food_nutrition_facts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_item_id    UUID NOT NULL REFERENCES food_items(id) ON DELETE CASCADE,
    reference_amount DOUBLE PRECISION NOT NULL,
    reference_unit  TEXT NOT NULL DEFAULT 'g',
    calories_kcal   DOUBLE PRECISION,
    protein_g       DOUBLE PRECISION,
    carbs_g         DOUBLE PRECISION,
    fat_g           DOUBLE PRECISION,
    fiber_g         DOUBLE PRECISION,
    sodium_mg       DOUBLE PRECISION,
    metadata_json   JSONB DEFAULT '{}'::jsonb,

    UNIQUE (food_item_id)
);

CREATE TABLE food_household_measures (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_item_id    UUID NOT NULL REFERENCES food_items(id) ON DELETE CASCADE,
    label           TEXT NOT NULL,
    grams_equivalent DOUBLE PRECISION,
    ml_equivalent   DOUBLE PRECISION,
    unit_count      DOUBLE PRECISION,
    sort_order      INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_food_measures_item ON food_household_measures(food_item_id);
```

- [ ] **Step 2: Rodar migration e commit**

```bash
make migrate-up
git add services/api/migrations/000031_*
git commit -m "feat(phase2): add migration for food catalog (items, nutrition, measures)"
```

---

## Task 13: Catalog module — domain, repository, usecase, handler

**Files:**
- Create: `services/api/internal/catalog/domain/food.go`
- Create: `services/api/internal/catalog/repository/postgres.go`
- Create: `services/api/internal/catalog/usecase/catalog.go`
- Create: `services/api/internal/catalog/handler.go`

- [ ] **Step 1: Domain**

- FoodItem struct — id, tenant_id (nullable), source_scope, name, common_name, brand, food_group, enabled
- NutritionFacts struct — reference_amount, reference_unit, calories, protein, carbs, fat, fiber, sodium
- HouseholdMeasure struct — label, grams_equivalent, ml_equivalent, unit_count, sort_order
- FoodItemWithDetails struct — FoodItem + NutritionFacts + []HouseholdMeasure (para respostas completas)
- SourceScope: Global, Tenant
- Sentinel errors: ErrFoodNotFound

- [ ] **Step 2: Repository**

- Search(ctx, tenantID, query, foodGroup, page, perPage) ([]FoodItemWithDetails, int, error) — busca global + tenant-scoped, ILIKE no name/common_name
- GetByID(ctx, tenantID, id) — retorna item com detalhes (global ou do tenant)
- Create(ctx, FoodItem, NutritionFacts, []HouseholdMeasure) — sempre tenant-scoped
- Update(ctx, FoodItem, NutritionFacts, []HouseholdMeasure) — apenas tenant-scoped
- Delete(ctx, tenantID, id) — apenas tenant-scoped (soft delete: enabled=false)
- ListGroups(ctx, tenantID) ([]string, error) — lista food_groups disponiveis

- [ ] **Step 3: Usecase**

- SearchFoods — busca com paginacao, combina global + tenant
- GetFood — por ID, valida acesso (global ou mesmo tenant)
- CreateTenantFood — cria alimento customizado do tenant, audita
- UpdateTenantFood — atualiza apenas foods do tenant, audita
- DeleteTenantFood — soft delete, audita
- ListFoodGroups

- [ ] **Step 4: Handler**

- GET    /foods?q=&group=&page=&per_page= — RequirePermission("diet:read") — busca global + tenant
- GET    /foods/{id} — RequirePermission("diet:read")
- POST   /foods — RequirePermission("clinical:write") — cria food tenant-scoped
- PUT    /foods/{id} — RequirePermission("clinical:write") — atualiza food tenant-scoped
- DELETE /foods/{id} — RequirePermission("clinical:write") — desativa food tenant-scoped
- GET    /foods/groups — RequirePermission("diet:read")

- [ ] **Step 5: Testes**

- usecase/catalog_test.go — mock repo, testar busca combinada, bloqueio de edicao em food global

- [ ] **Step 6: Wire-up em main.go**

- [ ] **Step 7: Commit**

```bash
git add services/api/internal/catalog/ services/api/cmd/api/main.go
git commit -m "feat(phase2): food catalog module — search, CRUD, nutrition facts, household measures"
```

---

## Task 14: Seed de permissoes da Fase 2

**Files:**
- Create: `services/api/migrations/000032_seed_phase2_permissions.up.sql`
- Create: `services/api/migrations/000032_seed_phase2_permissions.down.sql`

- [ ] **Step 1: Migration 000032 — novas permissoes e atribuicoes**

```sql
-- Novas permissoes para Fase 2
INSERT INTO permissions (id, code, application_scope, description) VALUES
  ('10000000-0000-0000-0000-000000000030', 'professionals:read',    'tenant', 'Ler profissionais'),
  ('10000000-0000-0000-0000-000000000031', 'professionals:write',   'tenant', 'Gerenciar profissionais'),
  ('10000000-0000-0000-0000-000000000032', 'availability:manage',   'tenant', 'Gerenciar disponibilidade'),
  ('10000000-0000-0000-0000-000000000033', 'blocks:manage',         'tenant', 'Gerenciar bloqueios de agenda'),
  ('10000000-0000-0000-0000-000000000034', 'measurements:read',     'tenant', 'Ler bioimpedancia'),
  ('10000000-0000-0000-0000-000000000035', 'measurements:write',    'tenant', 'Registrar bioimpedancia'),
  ('10000000-0000-0000-0000-000000000036', 'foods:read',            'tenant', 'Ler catalogo alimentar'),
  ('10000000-0000-0000-0000-000000000037', 'foods:write',           'tenant', 'Gerenciar alimentos do tenant'),
  ('10000000-0000-0000-0000-000000000038', 'invites:manage',        'tenant', 'Gerenciar convites de pacientes')
ON CONFLICT (id) DO NOTHING;

-- Owner ganha todas as novas permissoes (ja tem todas via regra existente, mas garantir)
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000001', id
FROM permissions WHERE application_scope = 'tenant' AND id >= '10000000-0000-0000-0000-000000000030'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Nutritionist ganha permissoes clinicas da Fase 2
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000002', id
FROM permissions WHERE code IN (
  'professionals:read', 'availability:manage', 'blocks:manage',
  'measurements:read', 'measurements:write',
  'foods:read', 'foods:write', 'invites:manage'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Receptionist ganha leitura
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000003', id
FROM permissions WHERE code IN ('professionals:read', 'foods:read', 'measurements:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Patient ganha leitura de medidas publicadas
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000004', id
FROM permissions WHERE code IN ('measurements:read', 'foods:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;
```

Down migration reverte as insercoes.

- [ ] **Step 2: Rodar migration e commit**

```bash
make migrate-up
git add services/api/migrations/000032_*
git commit -m "feat(phase2): seed Phase 2 permissions for all roles"
```

---

## Task 15: Testes de integracao e Docker smoke test

**Files:**
- Modify: `Makefile` (adicionar target test-phase2)

- [ ] **Step 1: Rodar todos os testes unitarios**

```bash
cd services/api && go test ./... -count=1 -v
```

Resultado esperado: `ok` em todos os packages (identity, tenancy, rbac, billing, professional, scheduling, patient, clinical, bioimpedance, catalog).

- [ ] **Step 2: Rebuild e restart Docker**

```bash
make dev-down && make dev-detached
```

Esperar containers healthy.

- [ ] **Step 3: Rodar migrations no Docker**

```bash
make migrate-up
```

Verificar: todas as 32 migrations aplicadas sem erro.

- [ ] **Step 4: Verificar /health**

```bash
curl -s http://localhost:8081/health | jq .
```

Resultado: `{"status":"healthy"}` HTTP 200.

- [ ] **Step 5: Verificar novos endpoints (sem auth = 401)**

```bash
curl -s -w "\nHTTP %{http_code}\n" http://localhost:8081/professionals
curl -s -w "\nHTTP %{http_code}\n" http://localhost:8081/appointments
curl -s -w "\nHTTP %{http_code}\n" http://localhost:8081/patients
curl -s -w "\nHTTP %{http_code}\n" http://localhost:8081/foods
```

Resultado: todos retornam HTTP 401.

- [ ] **Step 6: Commit**

```bash
git add .
git commit -m "feat(phase2): all modules integrated, Docker smoke test passing"
```

---

## Task 16: Push, PR e tag da Fase 2

- [ ] **Step 1: Rodar testes finais**

```bash
cd services/api && go test ./... -count=1
```

- [ ] **Step 2: Push branch**

```bash
git push -u origin feat/fase-2-operacao-clinica
```

- [ ] **Step 3: Criar PR**

```bash
gh pr create --title "feat: Phase 2 — Operacao Clinica" --body "## Summary
- Professional module (CRUD, addresses, service modes)
- Scheduling module (availability rules, blocks, appointments with conflict detection)
- Patient module (CRUD, invite/code flow, portal activation, profiles, consents)
- Clinical module (anamnesis draft/finalize, progress notes, attachments)
- Bioimpedance module (manual measurements, BMI calculation, history, publication)
- Food catalog module (global + tenant-scoped, search, nutrition facts, household measures)
- 17 new migrations (000016-000032)
- Phase 2 RBAC permissions seeded for all roles
- All endpoints tenant-scoped, audited, RBAC-protected

## Test plan
- [ ] \`go test ./... -count=1\` passes
- [ ] \`make setup\` + \`make migrate-up\` applies all 32 migrations
- [ ] /health returns 200
- [ ] All new endpoints return 401 without auth
- [ ] Professional CRUD works with valid JWT + tenant header
- [ ] Appointment conflict detection blocks double-booking
- [ ] Patient invite code generation and activation flow
- [ ] Anamnesis finalization is irreversible
- [ ] Bioimpedance BMI auto-calculation
- [ ] Food search returns global + tenant items"
```

- [ ] **Step 4: Merge PR**

```bash
gh pr merge --merge
```

- [ ] **Step 5: Tag da Fase 2**

```bash
git checkout main && git pull
git tag -a v0.2.0-fase2 -m "Phase 2 Clinical Operations: professionals, scheduling, patients, clinical, bioimpedance, food catalog"
git push origin v0.2.0-fase2
```

---

## Checklist de cobertura — Fase 2

| Requisito (docs/05-delivery/mvp-phases.md) | Task | Status |
|---------------------------------------------|------|--------|
| Profissionais e equipe | 1, 2 | ✓ |
| Enderecos de atendimento | 1, 2 | ✓ |
| Modalidades (presencial/online) | 1, 2 | ✓ |
| Disponibilidade e regras recorrentes | 3, 5 | ✓ |
| Bloqueios de agenda | 3, 5 | ✓ |
| Appointments com deteccao de conflitos | 4, 5 | ✓ |
| Audit events de agendamento | 4, 5 | ✓ |
| Pacientes CRUD | 6, 7 | ✓ |
| Convites/codigo e portal paciente | 6, 7 | ✓ |
| Access links e ativacao | 6, 7 | ✓ |
| Perfis nutricionais | 6, 7 | ✓ |
| Consentimentos | 6, 7 | ✓ |
| Anamnese (draft/finalize) | 8, 9 | ✓ |
| Evolucao/progress notes | 8, 9 | ✓ |
| Anexos clinicos | 8, 9 | ✓ |
| Publicacao ao paciente | 9, 11 | ✓ |
| Bioimpedancia manual | 10, 11 | ✓ |
| Historico e comparacao temporal | 11 | ✓ |
| Catalogo alimentar (global + tenant) | 12, 13 | ✓ |
| Dados nutricionais | 12, 13 | ✓ |
| Medidas caseiras | 12, 13 | ✓ |
| Permissoes RBAC Fase 2 | 14 | ✓ |
| Smoke test Docker | 15 | ✓ |
| PR + tag | 16 | ✓ |
