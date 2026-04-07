# Módulos e escopo consolidado

## 1. Regra de leitura
Este documento descreve os módulos que a Nutrometra **já deve considerar como parte do produto inicial**. Não interpretar nenhuma seção abaixo como backlog opcional.

## 2. Superfícies do produto
### 2.1 Portal profissional
Responsável por operação clínica e operacional do tenant.

Inclui obrigatoriamente:
- login e segurança;
- gestão do tenant;
- profissionais e equipe;
- endereços e horários;
- agenda presencial e online;
- pacientes;
- prontuário, anamnese e evolução;
- bioimpedância;
- dietas;
- alimentos e medidas caseiras;
- documentos clínicos;
- exportação PDF e impressão;
- assinatura e plano;
- notificações e integrações.

### 2.2 Portal paciente
Responsável pela experiência gratuita do paciente vinculada à nutricionista.

Inclui obrigatoriamente:
- ativação por código;
- login;
- cardápio por refeições;
- visualização de substituições e observações;
- visualização da nutricionista;
- visualização da evolução publicada;
- visualização de documentos publicados;
- experiência prioritária em celular.

### 2.3 Backoffice
Responsável pela operação interna da plataforma.

Inclui obrigatoriamente:
- suporte;
- financeiro;
- billing;
- vendas;
- gestão de planos e entitlements;
- visão operacional dos tenants;
- auditoria.

## 3. Módulos do domínio
### 3.1 Identidade e acesso
- contas de profissional, equipe, paciente e backoffice;
- recuperação de acesso;
- papéis e permissões;
- sessões e segurança.

### 3.2 Multi-tenant
- tenant individual, clínica e empresa;
- isolamento de dados;
- configuração por tenant;
- overrides por tenant.

### 3.3 Planos, trial e billing
- planos Free, Starter, Pro, Enterprise;
- trial por plano e/ou por feature;
- limites por feature;
- cobrança recorrente;
- invoices e pagamentos;
- upgrades, downgrades e cancelamento.

### 3.4 Profissionais, equipe e unidades de atendimento
- profissional principal e equipe;
- cadastro de endereços;
- modalidade online;
- horários por endereço e modalidade;
- bloqueios e indisponibilidade.

### 3.5 Agenda
- criação de atendimentos;
- alteração, exclusão, cancelamento e remarcação;
- conflitos e buffers;
- sincronização com Google Calendar;
- trilha de alterações.

### 3.6 Pacientes e vínculo
- cadastro de paciente;
- associação ao profissional;
- convite/código;
- portal paciente;
- status do vínculo.

### 3.7 Prontuário clínico
- anamnese;
- evolução;
- observações;
- anexos;
- publicação seletiva ao paciente.

### 3.8 Bioimpedância e evolução corporal
- inserção manual dos dados aferidos;
- histórico de medições;
- comparação temporal;
- indicadores corporais;
- origem da informação.

### 3.9 Catálogo alimentar
- alimentos;
- grupos;
- dados nutricionais;
- medidas caseiras;
- porções;
- sinônimos e observações de uso.

### 3.10 Dietas
- plano alimentar por refeições;
- itens e quantidades;
- substituições;
- observações por refeição;
- publicação ao paciente;
- versionamento.

### 3.11 Documentos clínicos
- solicitação de exames;
- prescrições e demais documentos cabíveis;
- geração de PDF;
- impressão;
- histórico de versões.

### 3.12 Backoffice
- suporte a incidentes;
- gestão de plano e entitlement;
- análise de cobrança e inadimplência;
- visão de vendas;
- visão operacional e auditoria.

### 3.13 IA assistiva
- apoio à geração de rascunhos clínicos;
- apoio à elaboração de dietas;
- apoio por perfil nutricional;
- revisão obrigatória por profissional humano.

## 4. Limites que podem variar por plano
- quantidade de profissionais ativos;
- quantidade de pacientes ativos;
- quantidade de endereços;
- quantidade de agendas sincronizadas;
- quantidade de exports em PDF por período;
- acesso a IA assistiva;
- retenção de histórico;
- acesso a dashboards;
- acesso a suporte prioritário;
- acesso a módulos avançados de backoffice interno quando aplicável.

## 5. Funcionalidades que podem nascer com implementação simplificada
Estas funcionalidades **continuam obrigatórias no MVP**, apenas com execução técnica inicial mais simples:
- billing com um provedor;
- Google Calendar como primeira integração de agenda;
- bioimpedância por cadastro manual;
- PDFs com um engine principal;
- notificações com poucos canais iniciais;
- dashboards básicos.
