# PRD — Nutrometra

## 1. Visão do produto
A Nutrometra é uma plataforma SaaS de nutrição para profissionais autônomos, clínicas e empresas de nutrição. O produto centraliza operação clínica, agenda, prontuário, bioimpedância, dietas, documentos, assinatura, billing, suporte e relacionamento com o paciente em uma única solução web responsiva.

O sistema possui três experiências nativas:
- **portal profissional** para nutricionistas e equipe autorizada;
- **portal paciente** gratuito, ativado por convite/código gerado pelo profissional;
- **backoffice** para suporte, vendas, billing, gestão financeira e operações da plataforma.

## 2. Regra de escopo
Tudo o que foi levantado até aqui faz parte do escopo base do produto. O projeto não deve ser descrito como se tivesse:
- núcleo mínimo sem backoffice;
- núcleo mínimo sem portal paciente;
- núcleo mínimo sem billing;
- núcleo mínimo sem agenda;
- núcleo mínimo sem bioimpedância;
- núcleo mínimo sem dietas e PDFs.

O MVP é completo em cobertura funcional e controlado em profundidade de implementação.

## 3. Objetivos de negócio
- reduzir tempo operacional do profissional e da clínica;
- centralizar prontuário, agenda, bioimpedância, cardápio e documentos;
- transformar o produto em um SaaS escalável por plano e assinatura;
- melhorar adesão do paciente com acesso mobile aos cardápios e evolução;
- permitir governança clínica, rastreabilidade, segurança e atendimento interno;
- dar suporte a múltiplos perfis de negócio, de profissional individual a empresa com múltiplos nutricionistas.

## 4. Usuários do sistema
### 4.1 Externos pagantes
- nutricionista autônomo;
- clínica de nutrição;
- empresa/grupo de nutrição com múltiplos profissionais;
- secretária/recepção do tenant;
- gestor financeiro ou operacional do tenant.

### 4.2 Externos não pagantes
- paciente convidado por profissional.

### 4.3 Internos
- suporte;
- financeiro;
- vendas/comercial;
- operações/backoffice;
- administradores da plataforma.

## 5. Problemas resolvidos
- fragmentação entre agenda, prontuário, planilhas, PDFs e comunicação com pacientes;
- baixa padronização de condutas e documentos entre profissionais;
- pouca rastreabilidade sobre evolução corporal e bioimpedância;
- dificuldade de controlar recursos por plano, trial e assinatura;
- baixa qualidade da experiência mobile para o paciente;
- pouco controle administrativo sobre cobrança, suporte, vendas e permissões;
- dificuldade de manter catálogo alimentar consistente com medidas caseiras e porções.

## 6. Escopo funcional do produto desde o início
### 6.1 Identidade, acesso e SaaS
- cadastro, login, recuperação de conta e verificação de e-mail;
- autenticação segura para profissionais, equipe, pacientes e usuários internos;
- tenants do tipo profissional, clínica ou empresa;
- planos comerciais com feature gating e limites;
- trial por plano e/ou por funcionalidade;
- overrides por tenant;
- assinatura, cobrança, histórico financeiro e eventos de billing.

### 6.2 Operação do profissional
- cadastro profissional e de equipe;
- cadastro de um ou mais endereços de atendimento;
- configuração de horários por endereço e por modalidade online;
- agenda com criação, alteração, exclusão, cancelamento, bloqueio e remarcação;
- sincronização da agenda do profissional com Google Calendar;
- cadastro de pacientes;
- convite do paciente por código;
- prontuário, anamnese, evolução e anexos;
- registro manual de bioimpedância;
- registro histórico de bioimpedância;
- registro de meditas corporais com histórico;
- composição corporal e evolução histórica;
- criação, edição, publicação e versionamento de dietas;
- catálogo de alimentos, porções e medidas caseiras;
- geraçao de PDF e impressão;
- solicitação/prescrição de exames e outros documentos clínicos.

### 6.3 Portal do paciente
- cadastro com código de acesso;
- visualização de cardápios por refeição;
- visualização de substituições e observações da refeição;
- visualização dos dados da nutricionista;
- visualização da evolução publicada;
- visualização de documentos liberados;
- experiência mobile-first.

### 6.4 Backoffice
- atendimento de suporte;
- visão do ciclo de assinatura e cobrança;
- gestão de planos, entitlements e recursos;
- visão operacional de tenants e usuários;
- gestão financeira e de vendas conforme permissão;
- auditoria e trilha operacional.

## 7. Requisitos funcionais consolidados
### 7.1 Identidade e conta
- RF-001: o sistema deve permitir criação de conta segura para profissionais e equipe.
- RF-002: o sistema deve permitir autenticação segura, recuperação de acesso e gestão de sessão.
- RF-003: o sistema deve permitir ativação de paciente apenas mediante convite/código válido.
- RF-004: o sistema deve permitir diferentes papéis por aplicação e por tenant.

### 7.2 Tenants, planos e billing
- RF-010: o sistema deve permitir criação de tenant individual ou organizacional.
- RF-011: o sistema deve permitir contratar plano e iniciar trial.
- RF-012: o sistema deve permitir habilitar/desabilitar features por plano.
- RF-013: o sistema deve permitir limites por plano, como quantidade de profissionais, pacientes, exports e integrações.
- RF-014: o sistema deve permitir overrides manuais por tenant para suporte comercial/operacional.
- RF-015: o sistema deve registrar o histórico de assinatura, cobrança, mudança de plano e inadimplência.

### 7.3 Profissionais, endereços e agenda
- RF-020: o profissional deve poder cadastrar um ou mais endereços de atendimento.
- RF-021: cada endereço deve possuir regras próprias de agenda.
- RF-022: o profissional deve poder definir disponibilidade online sem endereço físico.
- RF-023: o sistema deve permitir criar, editar, cancelar, excluir e remarcar atendimentos.
- RF-024: o sistema deve permitir bloqueios de agenda.
- RF-025: o sistema deve permitir sincronização com Google Calendar.
- RF-026: o sistema deve permitir cadastro de horários por profissional, endereço e modalidade.

### 7.4 Pacientes e prontuário
- RF-030: o sistema deve permitir cadastro de pacientes.
- RF-031: o sistema deve permitir anamnese, evolução, observações e anexos.
- RF-032: o sistema deve permitir registrar bioimpedância manualmente.
- RF-033: o sistema deve manter histórico de evolução corporal.
- RF-034: o sistema deve permitir classificação do paciente em um ou mais perfis nutricionais.
- RF-035: o sistema deve permitir publicação seletiva de conteúdos para o portal do paciente.

### 7.5 Catálogo alimentar e dietas
- RF-040: o sistema deve possuir catálogo de alimentos com dados nutricionais e metadados relevantes.
- RF-041: o sistema deve permitir medidas caseiras, porções e quantidades customizadas.
- RF-042: o sistema deve permitir montar dietas por refeições.
- RF-043: o sistema deve permitir substituições alimentares.
- RF-044: o sistema deve permitir observações por refeição e por item.
- RF-045: o sistema deve permitir versionamento, publicação, exportação e impressão de dietas.

### 7.6 Documentos clínicos
- RF-050: o sistema deve permitir gerar solicitações de exames.
- RF-051: o sistema deve permitir gerar outros documentos clínicos cabíveis.
- RF-052: o sistema deve permitir exportação PDF e impressão desses documentos.
- RF-053: o sistema deve manter histórico de versões e exports.

### 7.7 Portal do paciente
- RF-060: o paciente deve poder visualizar cardápios separados por refeições.
- RF-061: o paciente deve poder visualizar dados da nutricionista.
- RF-062: o paciente deve poder visualizar evolução publicada.
- RF-063: o paciente deve poder visualizar documentos publicados.

### 7.8 Backoffice
- RF-070: o backoffice deve permitir suporte operacional por papel.
- RF-071: o backoffice deve permitir visualizar assinatura, pagamento, plano e uso.
- RF-072: o backoffice deve permitir ações financeiras e comerciais conforme nível de acesso.
- RF-073: o backoffice deve permitir auditoria e rastreabilidade das ações internas.

## 8. Requisitos não funcionais prioritários
- segurança por padrão;
- multi-tenant com isolamento forte;
- performance para centenas de usuários simultâneos com evolução para milhares;
- observabilidade, auditoria e trilhas de suporte;
- mobile-first;
- disponibilidade operacional;
- filas para jobs pesados e integrações;
- cache seletivo no front-end e no back-end;
- acessibilidade e usabilidade.

## 9. Fora do escopo inicial
Fora do escopo inicial apenas o que não foi solicitado e não é necessário para a operação núcleo, por exemplo:
- teleconsulta nativa por vídeo;
- importação automática de equipamentos de bioimpedância específicos;
- app mobile nativo;
- marketplace público;
- BI avançado.
