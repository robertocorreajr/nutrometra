# Integrações

## 1. Princípio
Integrações são parte do produto inicial quando suportam requisitos já definidos. A implementação inicial pode ser enxuta, mas o contrato arquitetural deve nascer pronto para expansão.

## 2. Google Calendar
Uso por profissional, conforme plano e autorização.

### Capacidades mínimas
- conectar conta do Google;
- armazenar credenciais de forma segura;
- criar eventos sincronizados;
- atualizar eventos sincronizados;
- cancelar eventos sincronizados;
- registrar `external_event_id`;
- reconciliar falhas de sincronização.

### Regras técnicas
- integração isolada em adapter próprio;
- tokens criptografados;
- fila para sincronização e reconciliação;
- retries com backoff;
- idempotência por evento e operação;
- auditoria de vinculação e desvinculação.

## 3. Billing provider
O sistema deve começar com um provedor de billing, mas com contrato interno que permita trocar ou adicionar provedores.

### Capacidades mínimas
- criação de customer;
- criação e atualização de subscription;
- trial;
- upgrade e downgrade;
- cancelamento;
- recebimento de webhooks de pagamento e invoice;
- atualização do estado financeiro local.

### Regras técnicas
- processar webhooks com idempotência;
- usar outbox/eventos para refletir mudanças internas;
- tratar falhas em fila;
- manter histórico de invoice e pagamento no banco próprio.

## 4. E-mail e notificações
### Casos mínimos
- convite do paciente;
- recuperação de conta;
- notificações operacionais;
- eventos de cobrança;
- confirmação de exportação pronta quando aplicável.

### Regras técnicas
- templates versionados;
- fila de envio;
- reprocessamento de falhas;
- logs de entrega quando possível.

## 5. Geração de PDF
A geração de PDF é uma integração interna obrigatória do produto.

### Casos mínimos
- dieta;
- solicitação de exames;
- demais documentos clínicos cabíveis.

### Regras técnicas
- geração desacoplada do request síncrono quando o documento for pesado;
- template versionado;
- armazenamento do arquivo final em object storage;
- histórico de exports;
- suporte à impressão.

## 6. IA assistiva
A IA assistiva deve operar por contexto clínico controlado.

### Capacidades mínimas
- receber contexto sanitizado;
- operar por perfil nutricional;
- gerar rascunho de dieta, orientações ou estrutura de conduta;
- registrar prompt, contexto resumido e versão do agente quando aplicável.

### Regras técnicas
- revisão humana obrigatória;
- sem execução autônoma de conduta clínica final;
- logs e auditoria de uso;
- proteção contra vazamento de dados além do necessário.
