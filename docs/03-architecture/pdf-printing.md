# PDF e impressão

## Requisitos
- dieta exportável em PDF;
- solicitações/prescrições exportáveis em PDF;
- layout legível em desktop e impressão A4;
- preservação de histórico de exportações.

## Estratégia
- gerar HTML canônico por template;
- renderizar PDF em worker;
- armazenar arquivo em object storage;
- devolver status de processamento para a UI;
- registrar auditoria e metadados da exportação.

## Fontes de verdade
- documento clínico versionado;
- versão publicada da dieta;
- dados do profissional e paciente no momento da geração.
