# Arquitetura de sistema

## Visão geral
A plataforma é composta por aplicações web separadas por contexto e uma API principal em Go, apoiada por workers para tarefas assíncronas.

## Blocos
- `apps/web-professional`
- `apps/web-patient`
- `apps/backoffice`
- `services/api`
- `services/worker`
- PostgreSQL
- Redis
- object storage
- provedor de autenticação
- provedor de billing
- Google Calendar adapter

## Princípios
- separação entre leitura/escrita quando fizer sentido;
- outbox para eventos confiáveis;
- filas para PDF, notificações, IA, cobranças e sincronizações;
- adapters externos isolados;
- autorização centralizada.

## Fluxo de alto nível
1. frontend chama API;
2. API autentica, resolve tenant, valida plano e papel;
3. caso seja mutação, executa caso de uso de aplicação;
4. persiste dados e publica outbox;
5. worker consome eventos/jobs;
6. resultado volta por polling, webhook interno ou atualização de estado.
