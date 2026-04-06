# Fluxos clínicos e operacionais

## 1. Princípio
Os fluxos abaixo descrevem o funcionamento esperado do produto inicial. Todos fazem parte do desenho-base da plataforma.

## 2. Onboarding do profissional e do tenant
1. usuário cria conta;
2. cria tenant ou aceita convite para tenant existente;
3. escolhe plano ou inicia trial;
4. completa perfil profissional;
5. cadastra endereço(s) de atendimento e/ou modalidade online;
6. configura disponibilidade por endereço e modalidade;
7. conecta Google Calendar se desejar e se o plano permitir;
8. convida equipe quando aplicável.

## 3. Cadastro do paciente
1. profissional ou equipe autorizada cria paciente;
2. preenche dados básicos e vínculo principal;
3. registra consentimentos e observações iniciais;
4. classifica perfis nutricionais relevantes;
5. opcionalmente gera código de acesso para o portal paciente.

## 4. Gestão de agenda
1. recepção ou profissional cria atendimento;
2. escolhe modalidade presencial ou online;
3. escolhe endereço quando presencial;
4. sistema valida disponibilidade, buffers, conflitos e regra de plano;
5. sistema cria o agendamento;
6. opcionalmente sincroniza com Google Calendar;
7. qualquer alteração relevante gera auditoria.

### Operações obrigatórias da agenda
- inserir atendimento;
- alterar atendimento;
- excluir atendimento quando ainda não houver vínculo clínico consolidado;
- cancelar atendimento;
- remarcar atendimento;
- bloquear período;
- gerenciar disponibilidade recorrente e exceções.

## 5. Consulta e prontuário
1. abrir paciente;
2. revisar histórico clínico e atendimentos anteriores;
3. preencher ou atualizar anamnese;
4. registrar evolução;
5. anexar documentos quando necessário;
6. registrar bioimpedância;
7. salvar rascunho ou concluir registro.

## 6. Bioimpedância
A plataforma deve suportar, no mínimo, cadastro manual dos valores aferidos pelo equipamento.

Fluxo:
1. profissional informa data e contexto da medição;
2. insere os indicadores coletados;
3. sistema valida ranges básicos e consistência técnica;
4. sistema registra a origem como manual;
5. sistema atualiza histórico de evolução;
6. sistema disponibiliza comparação entre medições;
7. os dados publicados ao paciente seguem política de publicação definida pela nutricionista.

### Dados típicos esperados
- peso;
- altura quando aplicável ao contexto;
- IMC calculado ou armazenado;
- percentual de gordura;
- massa muscular;
- água corporal;
- gordura visceral;
- taxa metabólica basal quando informada;
- idade metabólica quando informada;
- observações clínicas.

## 7. Montagem da dieta
1. profissional abre contexto do paciente;
2. define objetivo nutricional;
3. escolhe estrutura de refeições;
4. utiliza catálogo de alimentos e medidas caseiras;
5. adiciona itens, quantidades, unidade e observações;
6. adiciona substituições quando necessário;
7. opcionalmente usa IA assistiva para rascunho;
8. revisa clinicamente;
9. versiona;
10. publica ao paciente;
11. exporta PDF e/ou imprime.

### Estrutura de refeições prevista
Exemplos que o sistema deve suportar:
- café da manhã;
- desjejum;
- colação;
- almoço;
- lanche;
- jantar;
- ceia.

A estrutura deve ser configurável por dieta, não fixa.

## 8. Solicitação de exames e documentos clínicos
1. profissional escolhe o tipo documental;
2. sistema carrega dados básicos do paciente e do profissional;
3. profissional complementa instruções e observações;
4. sistema gera versão do documento;
5. profissional exporta PDF e/ou imprime;
6. opcionalmente publica cópia ao paciente.

## 9. Ativação do paciente
1. profissional gera convite/código;
2. paciente acessa portal web;
3. cria credenciais;
4. informa código;
5. sistema valida tenant, expiração, uso, status e vínculo;
6. conta do paciente fica vinculada ao tenant e ao paciente clínico;
7. paciente passa a ver apenas conteúdo publicado para ele.

## 10. Experiência do paciente
O paciente deve conseguir:
- visualizar cardápios por refeição;
- ver itens, quantidades e observações;
- ver substituições quando liberadas;
- ver dados da nutricionista;
- ver evolução publicada;
- ver documentos publicados.

## 11. Fluxo de suporte no backoffice
1. usuário abre chamado ou time interno identifica problema;
2. suporte acessa tenant conforme permissão;
3. analisa plano, billing, erros, histórico e auditoria;
4. executa a ação permitida pelo papel;
5. registra tratativa;
6. ações sensíveis ficam auditadas.

## 12. Fluxo financeiro e comercial interno
1. backoffice identifica lead ou tenant;
2. avalia plano, trial, uso e cobrança;
3. realiza upgrade, ajuste, override ou acompanhamento financeiro conforme permissão;
4. sistema registra motivo, operador e timestamp;
5. mudanças impactam imediatamente os entitlements calculados.
