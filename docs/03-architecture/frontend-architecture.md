# Arquitetura de frontend

## Aplicações
### Web professional
- foco em produtividade;
- formulários extensos;
- agenda;
- prontuário;
- geração de dieta/documentos.

### Web patient
- foco em leitura simples e rápida;
- majoritariamente mobile;
- navegação curta;
- cardápio por refeição.

### Backoffice
- foco em tabelas, filtros, auditoria e gestão operacional.

## Diretrizes
- design system compartilhado quando possível;
- autenticação por aplicação;
- feature gating também no frontend, sem depender apenas dele;
- SSR/CSR conforme necessidade do stack escolhido;
- caching de leitura com invalidação explícita.
