# Estratégia de IA assistiva

## Objetivo
A IA deve acelerar o trabalho do nutricionista, nunca substituir decisão clínica.

## Casos de uso iniciais
- rascunho de dieta baseado em contexto do paciente;
- sugestão de estrutura de refeições;
- sugestão de substituições alimentares;
- resumo operacional do prontuário;
- checklist de revisão antes da publicação.

## Guardrails
- sempre mostrar que a saída é assistiva;
- nunca publicar diretamente sem revisão humana;
- registrar contexto utilizado;
- limitar escopo por perfil nutricional;
- evitar linguagem prescritiva absoluta quando faltarem dados.

## Contexto mínimo para geração
- perfil do paciente;
- objetivo;
- histórico relevante;
- restrições/alergias/intolerâncias;
- preferências e rotina;
- perfis de público aplicáveis.
