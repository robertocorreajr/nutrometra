# Extensão de contexto do projeto

## Decisões operacionais permanentes
- mobile-first para todas as telas;
- backend em Go;
- multi-tenant nativo;
- billing e planos fazem parte do núcleo;
- agenda, bioimpedância e dieta não são módulos opcionais;
- portal paciente é web e gratuito por convite;
- backoffice é obrigatório desde o início.

## Regras para o Claude Code
- não sugerir arquitetura monolítica sem modularização;
- não tratar billing, feature gating e RBAC como backlog futuro;
- não acoplar PDF ao request síncrono quando evitável;
- não permitir que IA grave/publice automaticamente sem revisão humana.
