# Filas, cache e observabilidade

## Filas
Usar filas para:
- geração de PDF;
- sincronização de calendário;
- envio de e-mails/notificações;
- cobrança e reconciliação;
- processamento de anexos;
- rotinas de IA assistiva;
- reprocessamentos e jobs agendados.

## Cache
Usar cache somente com chave e invalidação definidas.

### Bons candidatos
- entitlements por tenant;
- catálogos de alimentos;
- disponibilidade agregada da agenda;
- configurações públicas do tenant;
- sessões e rate limits.

## Observabilidade
- logs estruturados com `request_id`, `tenant_id`, `user_id`;
- métricas por endpoint, fila, integração e feature;
- tracing distribuído em API, worker e integrações;
- alertas para latência, erro, retries, fila parada e falha em billing.
