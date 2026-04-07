# Multi-tenant, billing e RBAC

## 1. Princípio
O desenho de autorização da Nutrometra depende sempre da combinação entre:
1. autenticação;
2. pertencimento ao contexto correto;
3. papel;
4. permissão;
5. entitlement do plano;
6. override ativo quando existir.

Nenhum módulo relevante deve ser implementado fora desse modelo.

## 2. Modelo de tenancy
### Tipos de tenant
- profissional individual;
- clínica;
- empresa de nutrição.

Todos os dados operacionais e clínicos devem carregar `tenant_id` e obedecer isolamento lógico forte.

## 3. Entitlements e limites
A autorização final para uso de uma funcionalidade depende da combinação de:
- identidade válida;
- papel do usuário;
- pertencimento ao tenant;
- entitlement ativo no plano;
- limite ainda disponível;
- override ativo, quando houver.

### Exemplos de entitlements
- quantidade máxima de profissionais;
- quantidade máxima de pacientes ativos;
- quantidade máxima de endereços;
- quantidade máxima de conexões Google Calendar;
- quantidade máxima de exports em PDF por período;
- acesso à IA assistiva;
- acesso a dashboards;
- acesso a suporte prioritário.

## 4. Papéis do tenant
- `tenant_owner`
- `tenant_admin`
- `nutritionist`
- `assistant`
- `receptionist`
- `financial_manager`
- `viewer`
- `patient`

## 5. Papéis do backoffice
- `backoffice_support`
- `backoffice_finance`
- `backoffice_sales`
- `backoffice_admin`

## 6. Matriz resumida por papel
### 6.1 Tenant owner
- gerencia tenant, plano, equipe e configuração;
- acessa visão operacional ampla do tenant;
- pode delegar funções;
- não ultrapassa limites do plano sem override.

### 6.2 Tenant admin
- gestão administrativa e operacional do tenant;
- pode cadastrar profissionais, endereços e agenda;
- acessa dados conforme política clínica do tenant.

### 6.3 Nutritionist
- acessa agenda, pacientes, prontuário, bioimpedância, dietas e documentos;
- publica conteúdo ao paciente;
- usa IA assistiva se o plano permitir.

### 6.4 Assistant / Receptionist
- agenda, cadastro operacional e suporte ao atendimento;
- acesso clínico restrito;
- sem capacidade de finalizar condutas clínicas sensíveis quando a permissão não existir.

### 6.5 Financial manager
- assinatura, invoices, pagamentos, upgrades e visão financeira do tenant;
- sem acesso desnecessário ao prontuário clínico.

### 6.6 Viewer
- visão limitada e somente leitura em escopos autorizados.

### 6.7 Patient
- acesso apenas ao próprio conteúdo publicado;
- sem visão lateral de outros pacientes, profissionais internos ou dados administrativos.

### 6.8 Backoffice support
- suporte operacional limitado e auditado;
- pode consultar tenant, usuários, erros e trilhas necessárias para atendimento.

### 6.9 Backoffice finance
- visão financeira, cobrança e inadimplência;
- sem acesso ampliado ao conteúdo clínico sem justificativa e permissão explícita.

### 6.10 Backoffice sales
- gestão comercial, trial, plano, oportunidade e upgrade;
- pode operar entitlements comerciais conforme política interna.

### 6.11 Backoffice admin
- acesso administrativo amplo e auditado;
- ações sensíveis exigem razão registrada.

## 7. Planos iniciais do produto
### Free
- 1 profissional;
- limite reduzido de pacientes;
- poucos endereços;
- sem Google Calendar ou com trial curto;
- IA indisponível ou limitada;
- PDFs limitados.

### Starter
- mais pacientes;
- até poucos profissionais;
- agenda completa;
- PDF liberado;
- alguns trials de integrações.

### Pro
- múltiplos profissionais;
- Google Calendar;
- IA assistiva;
- mais capacidade operacional;
- analytics básicos.

### Enterprise
- limites expandidos;
- maior granularidade;
- suporte avançado;
- flexibilidade comercial controlada.

## 8. Regras obrigatórias de segurança do backoffice
- acesso sempre auditado;
- menor privilégio possível;
- operações sensíveis com confirmação e justificativa;
- mascaramento parcial de dados quando o papel não exigir visualização integral;
- trilha de quem acessou, quando acessou e qual ação executou;
- possibilidade de revogação rápida de acesso.

## 9. Regras obrigatórias de billing
- mudança de plano deve refletir em entitlements calculados;
- trial deve ter início, fim e escopo explícitos;
- overrides devem ter razão, autor e validade;
- webhooks de cobrança devem ser idempotentes;
- inconsistências financeiras devem ir para reconciliação e fila de tratamento.
