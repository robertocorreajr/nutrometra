# Modelo de domínio nutricional

## Agregados principais
- Tenant
- Subscription
- Professional
- ProfessionalSchedule
- Patient
- ClinicalRecord
- BodyAssessment
- BioimpedanceMeasurement
- MealPlan
- ClinicalDocument
- PatientInvite

## Regras de domínio importantes
- todo recurso clínico pertence a um tenant;
- um paciente pode ter vários perfis nutricionais associados;
- bioimpedância é uma medição histórica, nunca sobrescrita;
- dieta publicada ao paciente deve ser versionada;
- documento clínico exportado gera evidência de exportação;
- agenda deve validar conflito temporal e modalidade;
- feature bloqueada por plano não pode ser executada, mesmo que a UI tente chamá-la.
