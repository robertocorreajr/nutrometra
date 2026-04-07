# Segurança

## Controles mínimos
- autenticação com provedor confiável;
- MFA preparado para ativação;
- controle de sessão seguro;
- RBAC + entitlements;
- rate limiting;
- auditoria;
- criptografia em trânsito;
- proteção de segredos;
- isolamento por tenant.

## Dados sensíveis
- prontuário e evolução clínica;
- bioimpedância;
- dados de contato do paciente;
- tokens de integração;
- billing e pagamentos.

## Regras
- não expor dados de um tenant para outro;
- não permitir bypass de feature gating pela UI;
- logs não devem conter segredos;
- acesso de backoffice deve ser auditado e justificável.
