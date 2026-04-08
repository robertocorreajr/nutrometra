# Portal Profissional v1 — Design Spec

## Contexto

O backend da Nutrometra está 100% completo (18+ módulos, 94 migrations, 60+ endpoints) e a infraestrutura frontend compartilhada (Turborepo, @nutrometra/ui, @nutrometra/api-client, @nutrometra/auth) está pronta. Os três apps frontend (profissional, paciente, backoffice) possuem apenas placeholders.

Esta spec define a primeira entrega funcional do **Portal Profissional** (`apps/web-professional`), cobrindo o shell de navegação, uma dashboard mínima, gestão de pacientes e prontuário clínico (anamnese, evoluções e anexos).

## Escopo

**Incluído:**
- Shell de navegação (sidebar expandida, header, mobile drawer)
- Dashboard mínima (KPIs, próximas consultas, pacientes recentes)
- Pacientes: lista, cadastro, edição, convites
- Prontuário: tabs (resumo, anamnese wizard, evoluções, anexos)
- Componentes reutilizáveis extraídos: DataTable, TabNav, FormWizard

**Excluído desta spec:**
- Agenda/scheduling (endpoints prontos, UI futura)
- Dietas (endpoints prontos, UI futura)
- Documentos clínicos (endpoints prontos, UI futura)
- Bioimpedância (endpoints prontos, UI futura)
- Exports/PDF (endpoints prontos, UI futura)
- IA assistiva (endpoints prontos, UI futura)
- Portal Paciente e Backoffice

## Abordagem

Vertical slices com extração incremental. Constrói feature por feature, extraindo componentes reutilizáveis (DataTable, FormWizard, TabNav) conforme a necessidade real aparece. Os shared packages (@nutrometra/ui, @nutrometra/api-client, @nutrometra/auth) são consumidos e estendidos organicamente.

---

## 1. Shell & Navegação

### 1.1 Estrutura de rotas (App Router)

```
apps/web-professional/src/
├── app/
│   ├── layout.tsx                    → Providers (Auth, API, QueryClient)
│   ├── (auth)/
│   │   └── signin/page.tsx           → Login via Zitadel
│   ├── (app)/
│   │   ├── layout.tsx                → AppShell (sidebar + header)
│   │   ├── page.tsx                  → Dashboard
│   │   ├── pacientes/
│   │   │   ├── page.tsx              → Lista de pacientes
│   │   │   ├── novo/page.tsx         → Cadastro de paciente
│   │   │   └── [id]/
│   │   │       ├── layout.tsx        → Header do paciente + TabNav
│   │   │       ├── page.tsx          → Tab Resumo
│   │   │       ├── anamnese/page.tsx → Tab Anamnese (wizard)
│   │   │       ├── evolucoes/page.tsx→ Tab Evoluções
│   │   │       └── anexos/page.tsx   → Tab Anexos
│   │   └── configuracoes/
│   │       └── page.tsx              → Perfil do profissional (placeholder)
├── components/                        → Componentes locais do app
│   ├── sidebar.tsx
│   ├── header.tsx
│   ├── patient-table.tsx
│   ├── patient-form.tsx
│   ├── anamnesis-wizard.tsx
│   ├── progress-note-form.tsx
│   └── attachment-grid.tsx
```

### 1.2 Sidebar

Usa o `AppShell` existente do @nutrometra/ui. A sidebar é um componente local (`sidebar.tsx`) passado como prop.

**Itens do menu:**

| Item | Ícone (lucide-react) | Rota | Estado |
|------|---------------------|------|--------|
| Dashboard | `LayoutDashboard` | `/` | ativo |
| Pacientes | `Users` | `/pacientes` | ativo |
| Agenda | `Calendar` | `/agenda` | desabilitado (opacity-50, sem link) |
| Dietas | `UtensilsCrossed` | `/dietas` | desabilitado |
| Documentos | `FileText` | `/documentos` | desabilitado |
| Configurações | `Settings` | `/configuracoes` | ativo |

Item ativo destacado com background `primary/10` e texto `primary`. Itens desabilitados são visíveis para comunicar completude futura do produto.

### 1.3 Header

- **Esquerda:** Nome do tenant (via `useTenant()`)
- **Direita:** Avatar/iniciais + nome do profissional + dropdown (Perfil, Sair)
- Usa `useAuth()` para dados do usuário

### 1.4 Mobile (< md / 768px)

- Sidebar esconde, hamburger icon aparece no header
- Toque no hamburger abre drawer lateral (overlay + sidebar completa)
- Drawer fecha ao clicar em item de menu ou no overlay
- Comportamento já suportado pelo `AppShell` existente

---

## 2. Dashboard

### 2.1 Layout

```
Desktop (md+):
┌──────────┬──────────┬──────────┐
│ Pacientes│ Consultas│ Pendentes│   ← 3 KPI cards
│ Ativos   │  Hoje    │          │
└──────────┴──────────┴──────────┘
┌─────────────────┬──────────────┐
│ Próximas        │ Pacientes    │   ← 2 colunas
│ Consultas       │ Recentes     │
└─────────────────┴──────────────┘

Mobile:
┌──────────┬──────────┐
│ Pacientes│ Consultas│   ← 2 colunas
├──────────┴──────────┤
│ Pendentes           │   ← 1 coluna
├─────────────────────┤
│ Próximas Consultas  │   ← stack
├─────────────────────┤
│ Pacientes Recentes  │
└─────────────────────┘
```

### 2.2 KPI Cards

Usam `Card` do @nutrometra/ui:

- **Pacientes ativos** — contagem total + link "Ver todos" → `/pacientes`
- **Consultas hoje** — contagem + horário da próxima (dados mockados nesta fase, hook preparado)
- **Pendências** — anamneses em draft + evoluções pendentes

### 2.3 Próximas Consultas

Lista dos próximos 5 agendamentos com: horário, nome do paciente, tipo (presencial/online). Dados mockados nesta spec (módulo de agenda não tem UI ainda). Hook `useUpcomingAppointments()` preparado para quando a UI de agenda for implementada.

### 2.4 Pacientes Recentes

Últimos 5 pacientes acessados/modificados. Avatar (iniciais), nome, última interação. Clique navega para `/pacientes/[id]`.

### 2.5 Estados

Cada card/seção implementa os 4 estados obrigatórios:
- **Loading** → `LoadingState` (skeleton)
- **Error** → `ErrorState` com botão retry
- **Empty** → `EmptyState` com ação (ex: "Nenhum paciente ainda" + "Cadastrar paciente")
- **Data** → conteúdo renderizado

---

## 3. Pacientes

### 3.1 Lista (`/pacientes`)

**Header:** `PageHeader` com título "Pacientes" e botão "+ Novo Paciente" como action.

**Tabela (desktop):**

| Coluna | Dados | Ordenável |
|--------|-------|-----------|
| Nome | Avatar (iniciais) + nome completo | sim |
| Status | Badge: Ativo (verde), Pendente (amarelo), Inativo (cinza) | sim |
| Última consulta | Data formatada ou "—" | sim |
| Perfil nutricional | Texto ou "—" | não |

**Funcionalidades:**
- Busca por nome — campo de texto, client-side para < 100 pacientes, server-side via `?q=` para mais
- Filtro por status — dropdown (Todos, Ativo, Pendente, Inativo)
- Ordenação por coluna — clique no cabeçalho
- Paginação — 20 por página, usa `PaginatedResponse<T>` do api-client
- Clique na linha — navega para `/pacientes/[id]`

**Mobile (< md):** Tabela se transforma em lista de cards. Cada card mostra avatar, nome, status badge, e última consulta.

**Componente extraído → `DataTable<T>`:**
- Props: `columns`, `data`, `isLoading`, `emptyState`, `searchPlaceholder`, `onSearch`, `pagination`, `onSort`
- Renderiza tabela no desktop, cards no mobile (via media query)
- Será reutilizado em evoluções, dietas, documentos, etc.

### 3.2 Cadastro (`/pacientes/novo`)

Formulário com React Hook Form + Zod validation.

**Campos:**

| Campo | Tipo | Obrigatório | Validação |
|-------|------|-------------|-----------|
| Nome completo | text | sim | min 3 chars |
| Email | email | sim | email válido |
| Telefone | tel | não | formato brasileiro |
| Data de nascimento | date | sim | no passado |
| Sexo | select (Masculino, Feminino, Outro) | sim | — |
| Perfil nutricional | select (13 perfis do CLAUDE.md) | não | — |
| Observações | textarea | não | max 1000 chars |

**Perfis nutricionais (select):** adulto geral, infantil, TEA e seletividade alimentar, gestantes, lactantes, idosos, atletas, fisiculturistas, emagrecimento e obesidade, diabetes e risco glicêmico, gastro, vegetariano e vegano, comportamento alimentar e risco clínico.

**Fluxo:**
1. Submete via `POST /patients`
2. Sucesso → toast de confirmação + redireciona para `/pacientes/[id]`
3. Erro de validação → mensagens inline por campo
4. Erro de servidor → `ErrorState` com mensagem

### 3.3 Edição

Mesmos campos do cadastro. Acessível via botão "Editar" no header do detalhe do paciente. Abre em modal ou página dedicada (modal para manter contexto). Pré-preenche com dados do `usePatient(id)`. Submete via `PUT /patients/{id}`.

### 3.4 Convites

Dentro da tab Resumo do paciente:
- **Paciente pendente:** botão "Enviar Convite" → `POST /patients/{id}/invites` → exibe código + link copiável (clipboard API)
- **Paciente ativo:** mostra data de ativação
- **Paciente inativo:** opção de reenviar convite

### 3.5 Hooks de API

Novos hooks a criar no @nutrometra/api-client:

```typescript
usePatients(params?: { q?: string; status?: string; limit?: number; offset?: number })
  → GET /patients

usePatient(id: string)
  → GET /patients/{id}

useCreatePatient()
  → POST /patients (mutation, invalidates usePatients)

useUpdatePatient(id: string)
  → PUT /patients/{id} (mutation, invalidates usePatient + usePatients)

useCreateInvite(patientId: string)
  → POST /patients/{patientId}/invites (mutation)

usePatientProfile(patientId: string)
  → GET /patients/{patientId}/profile
```

### 3.6 Tipos TypeScript

Novos tipos a adicionar no @nutrometra/api-client:

```typescript
type PatientStatus = "active" | "pending" | "inactive"

type Patient = {
  id: UUID
  tenant_id: UUID
  name: string
  email: string
  phone?: string
  date_of_birth: string
  sex: "M" | "F" | "O"
  nutritional_profile?: string
  notes?: string
  status: PatientStatus
  created_at: string
  updated_at: string
}

type PatientInvite = {
  id: UUID
  patient_id: UUID
  code: string
  link: string
  expires_at: string
  activated_at?: string
}

type PatientProfile = {
  id: UUID
  patient_id: UUID
  weight_kg?: number
  height_cm?: number
  bmi?: number
  body_fat_pct?: number
  notes?: string
}
```

---

## 4. Prontuário (Detalhe do Paciente)

### 4.1 Layout do paciente (`/pacientes/[id]/layout.tsx`)

**Header fixo:**
- Botão voltar "← Pacientes" (navega para `/pacientes`)
- Avatar (iniciais) + Nome + Idade (calculada) + Perfil nutricional + Badge de status
- Botão "Editar" (abre modal de edição)

**Tab bar (abaixo do header):**

| Tab | Rota | Descrição |
|-----|------|-----------|
| Resumo | `/pacientes/[id]` | Visão consolidada |
| Anamnese | `/pacientes/[id]/anamnese` | Wizard de anamnese |
| Evoluções | `/pacientes/[id]/evolucoes` | Notas de evolução |
| Anexos | `/pacientes/[id]/anexos` | Documentos e exames |

Tabs navegam via URL (App Router). Tab ativa baseada no `pathname`.

**Mobile:** Tab bar faz scroll horizontal se não couber na tela.

**Componente extraído → `TabNav`:**
- Props: `tabs: Array<{ label: string; href: string }>`
- Renderiza horizontal tabs com active state baseado em `usePathname()`
- Scroll horizontal no mobile
- Reutilizável em futuras áreas com tabs

### 4.2 Tab Resumo (`/pacientes/[id]`)

Cards informativos em grid (2 colunas desktop, 1 mobile):

| Card | Dados | Fonte |
|------|-------|-------|
| Dados Pessoais | Nome, idade, email, telefone, sexo | `usePatient(id)` |
| Última Anamnese | Data, status (draft/finalizada), link para tab | `useAnamneses(patientId)` |
| Bioimpedância | Último peso, IMC, % gordura (se disponível) | `usePatientProfile(patientId)` |
| Próxima Consulta | Data, horário, tipo (mockado nesta fase) | mock data |
| Convite/Acesso | Status + ação (enviar convite ou data de ativação) | `usePatient(id)` |

Cards sem dados mostram `EmptyState` inline com ação contextual (ex: "Nenhuma anamnese" + botão "Criar anamnese" que navega para a tab).

### 4.3 Tab Anamnese — Wizard (`/pacientes/[id]/anamnese`)

**Estados da página:**
- **Sem anamnese:** botão "Iniciar Anamnese" (cria via POST, redireciona para wizard)
- **Anamnese em draft:** retoma wizard do ponto onde parou
- **Anamnese finalizada:** exibe em modo leitura com botão "Nova Anamnese"

**Wizard — 5 etapas:**

**Etapa 1 — Dados Clínicos:**
| Campo | Tipo | Obrigatório |
|-------|------|-------------|
| Peso (kg) | number | sim |
| Altura (cm) | number | sim |
| Circunferência abdominal (cm) | number | não |
| Circunferência quadril (cm) | number | não |
| Pressão arterial | text (ex: 120/80) | não |
| Glicemia de jejum (mg/dL) | number | não |

**Etapa 2 — Histórico:**
| Campo | Tipo | Obrigatório |
|-------|------|-------------|
| Patologias | textarea (lista) | não |
| Cirurgias anteriores | textarea | não |
| Histórico familiar | textarea | não |
| Medicamentos em uso | textarea | não |

**Etapa 3 — Hábitos:**
| Campo | Tipo | Obrigatório |
|-------|------|-------------|
| Número de refeições/dia | number | sim |
| Hidratação (L/dia) | number | não |
| Atividade física | select (sedentário, leve, moderado, intenso) | sim |
| Frequência atividade (dias/semana) | number | não |
| Horas de sono | number | não |
| Tabagismo | boolean | não |
| Etilismo | select (não, social, regular) | não |
| Rotina alimentar | textarea | não |

**Etapa 4 — Alergias e Restrições:**
| Campo | Tipo | Obrigatório |
|-------|------|-------------|
| Alergias alimentares | textarea (lista) | não |
| Intolerâncias | textarea (lista) | não |
| Restrições (religiosas/éticas) | textarea | não |
| Preferências alimentares | textarea | não |
| Aversões alimentares | textarea | não |

**Etapa 5 — Revisão:**
- Resumo de todos os dados preenchidos (modo leitura)
- Campo "Observações gerais" (textarea)
- Campo "Objetivos do paciente" (textarea)
- Botão "Finalizar Anamnese" (chama `POST /anamneses/{id}/finalize`)

**Comportamento do wizard:**
- Progress bar no topo: 5 círculos/passos, preenchidos conforme avança
- Botões: "Anterior" | "Próximo" (passos 1-4), "Anterior" | "Finalizar" (passo 5)
- Autosave: salva rascunho automaticamente ao mudar de etapa via `PATCH`
- Validação Zod por etapa: campos obrigatórios impedem avanço, opcionais podem ficar vazios
- Mobile: mesmo wizard, layout empilhado, botões fixos no rodapé da tela

**Componente extraído → `FormWizard`:**
- Props: `steps: Array<{ title: string; schema: ZodSchema; component: ReactNode }>`, `onStepChange`, `onComplete`, `initialStep`
- Gerencia: navegação entre steps, progress bar, validação por step, callback de autosave
- Reutilizável para futuros formulários multi-step

### 4.4 Tab Evoluções (`/pacientes/[id]/evolucoes`)

**Listagem:**
- Lista cronológica (mais recente primeiro)
- Cada item: data, preview do texto (truncado ~150 chars), badge "Rascunho" ou "Finalizada"
- Botão "+ Nova Evolução" no `PageHeader`
- `EmptyState` quando não há evoluções

**Formulário de evolução (inline, expande abaixo do botão):**

| Campo | Tipo | Obrigatório |
|-------|------|-------------|
| Data da consulta | date (default: hoje) | sim |
| Observações clínicas | textarea | sim |
| Conduta | textarea | não |

**Ações:**
- "Salvar Rascunho" → `POST` ou `PATCH` (mantém editável)
- "Finalizar" → `POST /progress-notes/{id}/finalize` (trava edição, registra auditoria)
- Evolução finalizada: exibe em modo leitura, sem edição

**Hooks de anamnese:**
```typescript
useAnamneses(patientId: string)
  → GET /patients/{patientId}/anamneses

useCreateAnamnesis(patientId: string)
  → POST /patients/{patientId}/anamneses (mutation)

useUpdateAnamnesis(anamnesisId: string)
  → PATCH /anamneses/{anamnesisId} (mutation, autosave)

useFinalizeAnamnesis(anamnesisId: string)
  → POST /anamneses/{anamnesisId}/finalize (mutation)
```

**Hooks de evoluções:**
```typescript
useProgressNotes(patientId: string)
  → GET /patients/{patientId}/progress-notes

useCreateProgressNote(patientId: string)
  → POST /patients/{patientId}/progress-notes (mutation)

useUpdateProgressNote(noteId: string)
  → PATCH /progress-notes/{noteId} (mutation)

useFinalizeProgressNote(noteId: string)
  → POST /progress-notes/{noteId}/finalize (mutation)
```

### 4.5 Tab Anexos (`/pacientes/[id]/anexos`)

**Grid de anexos:**
- Cards em grid (3 colunas desktop, 2 mobile)
- Cada card: thumbnail (ícone por tipo de arquivo), nome, data upload, tamanho
- Botão "+ Upload" no `PageHeader`

**Upload:**
- Drag-and-drop zone + clique para selecionar
- Tipos aceitos: PDF, JPG, PNG
- Tamanho máximo: 10MB por arquivo
- Progress bar durante upload
- Submete via `POST /patients/{id}/attachments` (multipart/form-data)

**Ações por anexo:**
- Clique → abre em modal (imagem renderizada, PDF em iframe)
- Botão download
- Botão excluir com dialog de confirmação

**Hooks:**
```typescript
useAttachments(patientId: string)
  → GET /patients/{patientId}/attachments

useUploadAttachment(patientId: string)
  → POST /patients/{patientId}/attachments (mutation, multipart)

useDeleteAttachment(attachmentId: string)
  → DELETE /attachments/{attachmentId} (mutation)
```

**Tipos:**
```typescript
type Attachment = {
  id: UUID
  patient_id: UUID
  filename: string
  content_type: string
  size_bytes: number
  url: string
  created_at: string
}

type AnamnesisStatus = "draft" | "finalized"

type Anamnesis = {
  id: UUID
  patient_id: UUID
  status: AnamnesisStatus
  clinical_data?: Record<string, unknown>
  history?: Record<string, unknown>
  habits?: Record<string, unknown>
  allergies?: Record<string, unknown>
  observations?: string
  objectives?: string
  finalized_at?: string
  created_at: string
  updated_at: string
}

type ProgressNote = {
  id: UUID
  patient_id: UUID
  consultation_date: string
  clinical_observations: string
  conduct?: string
  status: "draft" | "finalized"
  finalized_at?: string
  created_at: string
  updated_at: string
}
```

---

## 5. Componentes Reutilizáveis Extraídos

Componentes criados nesta spec que serão adicionados ao @nutrometra/ui para reuso futuro:

| Componente | Localização | Propósito |
|------------|-------------|-----------|
| `DataTable<T>` | `packages/ui` | Tabela genérica com busca, ordenação, paginação, responsiva (cards no mobile) |
| `TabNav` | `packages/ui` | Navegação por tabs via URL, scroll horizontal mobile |
| `FormWizard` | `packages/ui` | Formulário multi-step com progress bar, validação por step, autosave |

---

## 6. Integração com Shared Packages

### @nutrometra/auth
- Root layout wraps com `SessionProvider`
- `(app)/layout.tsx` usa `useAuth()` para dados do profissional
- `middleware.ts` protege rotas `/` (redireciona para signin se não autenticado)
- `usePermission()` usado para feature flags futuras

### @nutrometra/api-client
- Root layout inicializa `configureApiClient({ baseUrl, getAccessToken, getTenantId })`
- `ApiProvider` wraps a árvore de componentes
- Novos hooks adicionados: patients, anamneses, progress-notes, attachments
- Novos tipos adicionados: Patient, PatientInvite, PatientProfile, Anamnesis, ProgressNote, Attachment

### @nutrometra/ui
- `AppShell` usado no `(app)/layout.tsx`
- `PageHeader`, `EmptyState`, `LoadingState`, `ErrorState` usados em todas as páginas
- `Button`, `Input`, `Card`, `Skeleton` usados nos formulários
- Novos componentes: `DataTable`, `TabNav`, `FormWizard`

---

## 7. Verificação

### Como testar end-to-end:

1. **Infraestrutura:** `make dev` sobe PostgreSQL, Redis, Zitadel, API
2. **Frontend:** `pnpm dev --filter web-professional` sobe o portal na porta 3000
3. **Login:** Acessar localhost:3000, redireciona para Zitadel, autenticar
4. **Dashboard:** Verificar KPIs (podem estar zerados), estados empty
5. **Pacientes:**
   - Clicar "+ Novo Paciente", preencher formulário, salvar
   - Verificar que aparece na lista
   - Clicar no paciente, verificar tabs
   - Enviar convite, verificar código gerado
6. **Anamnese:**
   - Clicar "Iniciar Anamnese"
   - Navegar pelo wizard preenchendo campos
   - Verificar autosave ao mudar de etapa
   - Finalizar e verificar modo leitura
7. **Evoluções:**
   - Criar nova evolução
   - Salvar rascunho
   - Finalizar e verificar trava de edição
8. **Anexos:**
   - Upload de arquivo (PDF ou imagem)
   - Verificar visualização em modal
   - Download e exclusão
9. **Mobile:** Testar em viewport 375px — drawer, cards responsivos, wizard mobile
