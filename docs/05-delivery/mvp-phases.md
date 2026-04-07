# Fases do MVP

## 1. Definição correta de MVP
O MVP da Nutrometra deve ser entendido como **produto funcional completo em cobertura**, e não como um recorte que elimina módulos centrais.

Isso significa que o MVP já inclui:
- portal profissional;
- portal paciente;
- backoffice;
- planos e billing;
- agenda;
- prontuário;
- bioimpedância;
- dietas;
- documentos clínicos;
- catálogo alimentar;
- PDFs e impressão;
- convites/código;
- RBAC, auditoria, filas e segurança.

## 2. O que pode ser simplificado no MVP
### Permitido simplificar
- um único provedor de billing;
- um único provedor de calendário;
- bioimpedância manual;
- poucos templates iniciais de PDF;
- dashboards básicos;
- poucos canais de notificação.

### Não permitido remover
- portal paciente;
- backoffice;
- planos e entitlements;
- agenda presencial e online;
- dietas por refeições;
- documentos clínicos;
- exportação PDF;
- gestão de pacientes por convite/código;
- catálogo de alimentos e medidas caseiras;
- auditoria e segurança.

## 3. Fase 1 — Fundação obrigatória
- autenticação e autorização;
- multi-tenant;
- planos, entitlements e assinatura;
- base do portal profissional;
- base do portal paciente;
- base do backoffice;
- observabilidade, auditoria e segurança.

## 4. Fase 2 — Operação clínica obrigatória
- profissionais, equipe e endereços;
- disponibilidade e agenda;
- pacientes e vínculo;
- prontuário, anamnese e evolução;
- bioimpedância;
- catálogo alimentar.

## 5. Fase 3 — Prescrição e publicação obrigatória
- dietas por refeições;
- substituições;
- publicação ao paciente;
- documentos clínicos;
- geração de PDF e impressão.

## 6. Fase 4 — Operação SaaS obrigatória
- billing provider;
- webhooks de cobrança;
- upgrades, downgrades e trial;
- backoffice financeiro, suporte e vendas;
- overrides por tenant.

## 7. Fase 5 — Integrações e IA obrigatórias do desenho inicial
- Google Calendar;
- IA assistiva por perfil nutricional;
- filas de reconciliação;
- refinamento de cache e performance.
