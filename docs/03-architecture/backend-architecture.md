# Arquitetura de backend (Go)

## Organização sugerida
- `internal/identity`
- `internal/tenancy`
- `internal/billing`
- `internal/catalog`
- `internal/scheduling`
- `internal/patients`
- `internal/clinical`
- `internal/documents`
- `internal/backoffice`
- `internal/ai`
- `internal/platform`

Cada módulo deve conter:
- domínio;
- casos de uso;
- repository interfaces;
- adapters.

## Padrões
- services de aplicação pequenos;
- repositories por agregado;
- transações conduzidas na aplicação;
- DTOs de entrada/saída separados do domínio;
- middlewares para autenticação, tenant resolution, autorização e observabilidade.
