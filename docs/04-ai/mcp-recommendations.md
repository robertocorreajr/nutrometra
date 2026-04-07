# MCPs e alternativas recomendadas

## MCPs prioritários
- filesystem: leitura e edição do repositório;
- git: inspeção de histórico e diffs;
- postgres: exploração de schema e validação de queries;
- playwright: testes E2E e validação de fluxos web;
- openapi: leitura guiada dos contratos da API.

## MCPs úteis adicionais
- redis: útil se houver servidor MCP confiável no ambiente para inspeção de cache/filas;
- task/issue tracker: para alinhar roadmap e execução;
- observability/logs: para depurar workers e integrações;
- object storage: quando houver necessidade de inspecionar PDFs e anexos.

## Alternativas que agregam mesmo sem MCP dedicado
- Makefiles/scripts padronizados para bootstrap e checks;
- fixtures e seeds para ambiente local;
- exemplos de payloads em `docs/examples/`;
- snapshots de OpenAPI e migrations sempre versionadas.

## Regra prática
Sem acesso a fontes de verdade, o Claude Code alucina estrutura. Por isso, priorizar:
1. `CLAUDE.md` forte;
2. `docs/` completos;
3. OpenAPI atualizada;
4. schema SQL/migrations;
5. seeds de dados;
6. subagentes claros.
