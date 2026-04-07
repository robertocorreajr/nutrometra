# CLAUDE.md

## Missão
Construir a Nutrometra como uma plataforma SaaS multi-tenant de nutrição para profissionais autônomos, clínicas e empresas de nutrição, com três superfícies nativas do produto:
1. portal profissional;
2. portal paciente;
3. backoffice administrativo.

A plataforma já nasce contemplando, como núcleo do produto:
- cadastro seguro, autenticação e autorização;
- planos, assinatura, trial, entitlements e billing;
- operação clínica e agenda;
- prontuário nutricional e evolução;
- bioimpedância e composição corporal;
- catálogo alimentar, medidas caseiras e montagem de dieta;
- exportação PDF e impressão;
- gestão do paciente por convite/código;
- suporte operacional, financeiro e vendas no backoffice;
- IA assistiva por perfis nutricionais.

## Regra central de escopo
Nada do que foi definido até aqui deve ser tratado como “extra”, “fase futura opcional” ou “complemento eventual”.

O MVP da Nutrometra é **completo em cobertura funcional**. A simplificação permitida é apenas de profundidade técnica inicial, por exemplo:
- integração com um único provedor antes de suportar vários;
- bioimpedância manual antes de importação automática;
- dashboard simples antes de BI avançado.

Não é permitido remover do desenho inicial:
- portal paciente;
- backoffice;
- billing e planos;
- agenda;
- bioimpedância;
- dietas;
- documentos clínicos;
- PDFs;
- catálogo de alimentos;
- convites/códigos;
- RBAC e auditoria.

## Restrições de negócio e produto
- O backend principal é **Go**.
- O produto é **web responsivo**, com prioridade **mobile-first**.
- O portal paciente é gratuito, porém só pode ser ativado por **convite/código emitido por profissional**.
- O sistema é **multi-tenant** desde o primeiro commit.
- Planos como Free, Starter, Pro e Enterprise fazem parte do núcleo do produto.
- Toda feature relevante deve considerar **trial**, **entitlements**, **limites por plano** e **overrides por tenant**.
- O produto deve suportar centenas de usuários simultâneos, com desenho preparado para milhares.
- Toda área clínica deve ser desenhada para revisão humana por profissional habilitado.

## Regras obrigatórias para gerar código
1. Antes de codar, sempre validar impacto em tenant, papel, auditoria, billing, cache e filas.
2. Nunca misturar regras clínicas com componentes de UI.
3. Separar claramente camadas de:
   - domínio;
   - aplicação;
   - adapters/integrations;
   - delivery (HTTP, filas, jobs, webhooks).
4. Tratar exportações PDF, impressão, sincronizações externas, notificações, cobrança, reconciliação, IA e tarefas pesadas como operações assíncronas quando isso trouxer robustez.
5. Toda operação mutável crítica deve registrar auditoria.
6. Toda integração externa deve ter retries, idempotência, reconciliação e DLQ quando necessário.
7. Toda consulta clínica, operacional e financeira deve respeitar isolamento de tenant e escopo do papel.
8. Toda sugestão de IA deve deixar claro:
   - contexto utilizado;
   - limitações;
   - necessidade de revisão humana.
9. O Claude Code deve preferir soluções simples, robustas, observáveis e auditáveis.
10. Quando houver dúvida entre flexibilidade e rastreabilidade, escolher rastreabilidade.
11. Sempre considerar que o paciente acessará majoritariamente pelo celular.
12. Sempre considerar que o nutricionista pode atender em múltiplos endereços e também online.
13. Sempre considerar que o sistema precisa lidar com assinatura, cobrança, trial e recursos por plano.
14. Ao concluir uma tarefa implementável, o fluxo padrão é sempre:
   - criar ou atualizar uma branch dedicada;
   - fazer um commit atômico das mudanças;
   - dar push da branch para o remoto;
   - só então iniciar a próxima tarefa.
   Se o push estiver bloqueado por falta de acesso, problema de rede ou política do repositório, interromper o fluxo e informar o bloqueio.

## Módulos obrigatórios do produto
- identidade e acesso;
- tenants, equipes, unidades e endereços de atendimento;
- assinatura, planos, trials, cupons, billing e entitlements;
- cadastro de profissionais e equipe operacional;
- cadastro de pacientes, convites, códigos e vínculo com nutricionista;
- disponibilidade, agenda, bloqueios, remarcações e calendário;
- agenda presencial e online;
- sincronização com Google Calendar;
- prontuário, anamnese, evolução e anexos;
- bioimpedância e composição corporal;
- catálogo de alimentos e medidas caseiras;
- dietas, refeições, substituições, observações e publicação ao paciente;
- documentos clínicos, solicitações, prescrições e exames;
- PDFs, impressão e histórico de exportação;
- notificações e comunicações;
- backoffice com suporte, financeiro, billing, vendas e governança;
- IA assistiva e governança de prompts;
- auditoria, observabilidade, filas, cache e segurança.

## Perfis nutricionais prioritários para IA assistiva
- adulto geral;
- infantil;
- TEA e seletividade alimentar;
- gestantes;
- lactantes;
- idosos;
- atletas;
- fisiculturistas;
- emagrecimento e obesidade;
- diabetes e risco glicêmico;
- gastro;
- vegetariano e vegano;
- comportamento alimentar e risco clínico.

## Decisões arquiteturais padrão
- monorepo no início;
- PostgreSQL como banco principal;
- Redis para cache, rate limiting, locks e filas leves quando necessário;
- filas para jobs críticos e demorados;
- object storage para anexos, PDFs e assets;
- OpenAPI como contrato canônico da API;
- eventos de domínio e outbox para integrações;
- feature flags e entitlements centralizados;
- RBAC por tenant e por aplicação;
- geração de PDF em serviço isolado ou worker dedicado;
- webhooks processados com idempotência;
- front-end com estratégia explícita de cache, invalidação e fallback offline seletivo.

## Convenções de implementação
### Backend Go
- usar `context.Context` em toda operação I/O;
- propagar `request_id`, `tenant_id`, `user_id`, `trace_id` e `actor_role`;
- validação de entrada na borda;
- regras de domínio fora de handlers;
- migrations versionadas;
- transações explícitas em operações clínicas, billing e agenda;
- testes unitários no domínio e integração nos adapters;
- erros internos nunca expõem segredos;
- jobs precisam ser idempotentes.

### Frontend
- mobile-first e acessibilidade AA como meta;
- estados de loading, error, empty e stale obrigatórios;
- autosave em formulários longos do profissional quando fizer sentido;
- experiência do paciente otimizada para celular;
- cache client-side apenas com estratégia clara de invalidação;
- impressão e exportação acessíveis a partir dos fluxos principais.

### Segurança
- autenticação com provedor compatível com OIDC/OAuth2;
- sessões seguras, rotação de tokens e proteção CSRF quando aplicável;
- MFA pronto para ativação;
- RBAC deny-by-default;
- trilha de auditoria para prontuário, agenda, documentos, exports e billing;
- segredos fora do repositório;
- criptografia para dados sensíveis e tokens de integração;
- acesso do backoffice sempre auditado.

## Saídas esperadas do Claude Code
Sempre que gerar uma implementação, trazer:
1. contexto e suposições;
2. plano de arquivos;
3. modelo de dados impactado;
4. endpoints, eventos e contratos impactados;
5. regras de autorização e entitlement impactadas;
6. estratégia de filas, cache e auditoria quando aplicável;
7. código;
8. testes;
9. riscos e pontos pendentes.

## Política de comandos sugerida
Use esta política para reduzir pedidos de autorização no Claude Code sem liberar comandos perigosos por padrão.

### Allowlist recomendada
- `pwd`
- `ls`
- `find`
- `rg`
- `sed`
- `cat`
- `head`
- `tail`
- `wc`
- `git status`
- `git diff`
- `git log`
- `git show`
- `git branch`
- `git remote -v`
- `git fetch`
- `git add`
- `git commit`
- `git push`
- comandos de teste e validação do projeto, como `go test`, `go test ./...`, `npm test`, `npm run test`, `make test`
- comandos de lint e build do projeto, como `go build`, `make build`, `npm run lint`, `npm run build`

### Denylist recomendada
- `rm -rf`
- `rm -fr`
- `git reset --hard`
- `git clean -fdx`
- `git checkout --`
- `git restore --source`
- `sudo`
- `chmod -R 777`
- `chown -R`
- `curl | sh`
- `wget | sh`
- `dd`
- `mkfs`
- comandos que instalem dependências globais, como `npm install -g`, `pnpm add -g`, `yarn global add`, `pip install --user`
- comandos que alterem produção, deploy ou infraestrutura sem aprovação explícita

### Regra operacional
Se um comando não estiver na allowlist e não for claramente de leitura, o Claude Code deve parar e pedir autorização antes de executar.
