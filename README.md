# Nutrometra Docs Pack

Este pacote define o escopo inicial **completo** da Nutrometra para uso no Claude Code.

## Princípios deste pacote
- o produto nasce **integral**, não como um conjunto de extras futuros;
- tudo o que foi solicitado até agora já deve ser tratado como **escopo base do projeto**;
- o MVP é **completo em cobertura funcional**, com simplificações apenas de profundidade técnica, nunca de ausência de módulo;
- o sistema é SaaS multi-tenant desde o início, com portal profissional, portal paciente e backoffice;
- o domínio clínico, operacional, financeiro, agenda, billing, PDFs, alimentos, bioimpedância, convites, permissões e integrações já fazem parte do desenho inicial.

## Ordem de leitura recomendada
1. `CLAUDE.md`
2. `docs/01-product/prd.md`
3. `docs/01-product/modules-and-scope.md`
4. `docs/02-domain/clinical-workflows.md`
5. `docs/03-architecture/system-architecture.md`
6. `docs/03-architecture/data-model.md`
7. `docs/03-architecture/multi-tenant-billing-rbac.md`
8. `docs/03-architecture/security.md`
9. `docs/04-ai/subagents.md`
10. `docs/05-delivery/mvp-phases.md`

## Regra para Claude Code
Ao gerar código, migrations, contratos, filas, testes e integrações, tratar os itens abaixo como obrigatórios desde o primeiro desenho:
- multi-tenant;
- planos, trials, entitlements e billing;
- profissionais, equipe e unidades/endereços;
- agenda presencial e online;
- paciente por convite/código;
- prontuário, anamnese e evolução;
- bioimpedância e composição corporal;
- catálogo de alimentos, porções e medidas caseiras;
- dietas por refeições, substituições e publicação ao paciente;
- PDFs e impressão de dieta e documentos clínicos;
- backoffice com níveis de acesso;
- auditoria, observabilidade, filas e segurança;
- IA assistiva por perfil nutricional.
