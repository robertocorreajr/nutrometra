# Requisitos não funcionais detalhados

## Performance
- operações síncronas de leitura comum devem ser rápidas mesmo em rede móvel;
- operações pesadas devem ir para fila;
- listagens críticas precisam de paginação e filtros.

## Escalabilidade
- desenho para centenas de usuários simultâneos no início;
- capacidade de crescer horizontalmente em API, workers e frontends estáticos.

## Segurança
- autenticação forte;
- autorização por papel e entitlements;
- isolamento por tenant;
- auditoria de ações críticas;
- segredos e chaves fora do código.

## Confiabilidade
- retries com backoff para integrações;
- DLQ para jobs falhos recorrentes;
- idempotência em pagamentos, webhooks e sincronizações.

## Observabilidade
- logs estruturados;
- métricas técnicas e de negócio;
- tracing distribuído;
- alertas para filas, erros e latência.

## UX
- mobile-first;
- acessibilidade como meta de produto;
- UX de leitura rápida no portal paciente;
- formulários extensos com recuperação de estado no portal profissional.
