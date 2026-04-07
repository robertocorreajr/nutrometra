# Modelo de dados consolidado

## 1. Diretrizes
- todo dado funcional relevante deve carregar `tenant_id`, exceto estruturas globais e contas centrais onde isso não se aplica;
- tabelas clínicas, financeiras e operacionais críticas devem ter `created_at`, `updated_at` e campos de ator quando necessário;
- histórico, auditoria e versionamento não são opcionais para domínios críticos;
- o modelo abaixo cobre o escopo base do produto.

## 2. Núcleo SaaS
### tenants
- id
- type (`solo_professional`, `clinic`, `company`)
- legal_name
- display_name
- slug
- status
- timezone
- locale
- trial_ends_at nullable
- created_at
- updated_at

### users
- id
- email
- email_verified_at nullable
- password_hash nullable
- external_auth_id nullable
- status
- last_login_at nullable
- created_at
- updated_at

### tenant_users
- id
- tenant_id
- user_id
- status
- joined_at
- invited_by_user_id nullable
- created_at
- updated_at

### roles
- id
- code
- application_scope (`tenant`, `backoffice`)
- name
- description

### permissions
- id
- code
- application_scope
- description

### role_permissions
- id
- role_id
- permission_id

### tenant_user_roles
- id
- tenant_id
- tenant_user_id
- role_id

### subscription_plans
- id
- code
- name
- active
- billing_cycle
- currency
- price_cents
- metadata_json
- created_at
- updated_at

### plan_features
- id
- plan_id
- feature_key
- enabled
- limit_value nullable
- trial_enabled
- trial_days nullable
- metadata_json

### coupons
- id
- code
- discount_type
- discount_value
- duration_type
- duration_cycles nullable
- active
- starts_at nullable
- ends_at nullable

### tenant_subscriptions
- id
- tenant_id
- plan_id
- coupon_id nullable
- status
- started_at
- trial_ends_at nullable
- renews_at nullable
- canceled_at nullable
- provider_customer_id nullable
- provider_subscription_id nullable
- created_at
- updated_at

### billing_invoices
- id
- tenant_id
- tenant_subscription_id
- provider_invoice_id nullable
- amount_cents
- currency
- status
- due_at nullable
- paid_at nullable
- hosted_url nullable
- created_at
- updated_at

### billing_payments
- id
- tenant_id
- billing_invoice_id
- provider_payment_id nullable
- amount_cents
- currency
- status
- paid_at nullable
- failure_reason nullable
- created_at

### tenant_feature_overrides
- id
- tenant_id
- feature_key
- enabled nullable
- limit_value nullable
- starts_at nullable
- ends_at nullable
- reason
- created_by_user_id nullable
- created_at

## 3. Profissionais, equipe e agenda
### professionals
- id
- tenant_id
- user_id nullable
- full_name
- registration_type (`CRN`, `CRO`, `CRM`, `other`) — para nutricionistas, sempre `CRN`
- registration_number — número do registro no conselho (ex: `12345/P`)
- registration_state — sigla do estado do conselho regional (ex: `SP`, `RJ`) — obrigatório para `CRN`
- document_number nullable
- phone nullable
- email nullable
- specialty_tags_json
- is_active
- created_at
- updated_at

### professional_addresses
- id
- tenant_id
- professional_id
- label
- street
- number
- complement nullable
- district
- city
- state
- zip_code
- country
- latitude nullable
- longitude nullable
- is_active
- created_at
- updated_at

### professional_service_modes
- id
- tenant_id
- professional_id
- mode (`onsite`, `online`)
- is_enabled
- created_at
- updated_at

### availability_rules
- id
- tenant_id
- professional_id
- address_id nullable
- service_mode
- weekday
- start_time
- end_time
- slot_minutes
- buffer_before_minutes
- buffer_after_minutes
- recurrence_end_date nullable
- is_active
- created_at
- updated_at

### schedule_blocks
- id
- tenant_id
- professional_id
- address_id nullable
- service_mode
- starts_at
- ends_at
- reason
- created_by_user_id
- created_at

### calendar_connections
- id
- tenant_id
- professional_id
- provider (`google`)
- status
- external_account_id
- encrypted_refresh_token
- token_expires_at nullable
- access_scope_json
- synced_at nullable
- created_at
- updated_at

### appointments
- id
- tenant_id
- professional_id
- patient_id
- address_id nullable
- service_mode (`onsite`, `online`)
- source (`manual`, `patient_request`, `import`, `sync`)
- status (`scheduled`, `confirmed`, `completed`, `cancelled`, `no_show`)
- scheduled_start_at
- scheduled_end_at
- cancellation_reason nullable
- notes nullable
- external_calendar_event_id nullable
- created_by_user_id
- updated_by_user_id nullable
- created_at
- updated_at

### appointment_audit_events
- id
- tenant_id
- appointment_id
- event_type
- actor_user_id nullable
- payload_json
- created_at

## 4. Pacientes e prontuário
### patients
- id
- tenant_id
- primary_professional_id
- external_code nullable
- full_name
- preferred_name nullable
- birth_date nullable
- sex nullable
- phone nullable
- email nullable
- contact_json
- status
- created_at
- updated_at

### patient_consents
- id
- tenant_id
- patient_id
- consent_type
- status
- accepted_at nullable
- revoked_at nullable
- metadata_json

### patient_invites
- id
- tenant_id
- patient_id nullable
- professional_id
- invite_code_hash
- expires_at
- max_uses
- used_count
- created_by_user_id
- created_at

### patient_access_links
- id
- tenant_id
- patient_id
- user_id
- activated_at
- status
- created_at

### patient_profiles
- id
- tenant_id
- patient_id
- profile_type
- priority
- metadata_json
- created_at

### anamneses
- id
- tenant_id
- patient_id
- professional_id
- appointment_id nullable
- status (`draft`, `finalized`)
- form_version
- payload_json
- created_at
- updated_at

### progress_notes
- id
- tenant_id
- patient_id
- professional_id
- appointment_id nullable
- note_type (`evolution`, `observation`, `followup`)
- content_json
- visibility_to_patient (`hidden`, `published`)
- created_at
- updated_at

### clinical_attachments
- id
- tenant_id
- patient_id
- professional_id
- appointment_id nullable
- file_key
- file_name
- mime_type
- size_bytes
- category
- created_at

## 5. Bioimpedância e evolução corporal
### body_measurements
- id
- tenant_id
- patient_id
- professional_id
- appointment_id nullable
- measured_at
- source (`manual`, `imported`)
- weight_kg nullable
- height_cm nullable
- bmi nullable
- body_fat_percent nullable
- skeletal_muscle_percent nullable
- muscle_mass_kg nullable
- fat_mass_kg nullable
- visceral_fat_index nullable
- body_water_percent nullable
- basal_metabolic_rate_kcal nullable
- metabolic_age nullable
- waist_cm nullable
- hip_cm nullable
- notes nullable
- raw_payload_json nullable
- created_at

### body_measurement_publications
- id
- tenant_id
- body_measurement_id
- patient_id
- published_by_user_id
- published_at
- visibility_mode (`summary`, `full`)

## 6. Catálogo de alimentos e dieta
### food_items
- id
- tenant_id nullable
- source_scope (`global`, `tenant`)
- name
- common_name nullable
- brand nullable
- food_group nullable
- enabled
- metadata_json
- created_at
- updated_at

### food_nutrition_facts
- id
- food_item_id
- reference_amount
- reference_unit
- calories_kcal nullable
- protein_g nullable
- carbs_g nullable
- fat_g nullable
- fiber_g nullable
- sodium_mg nullable
- metadata_json

### food_household_measures
- id
- food_item_id
- label
- grams_equivalent nullable
- ml_equivalent nullable
- unit_count nullable
- sort_order

### diets
- id
- tenant_id
- patient_id
- professional_id
- title
- objective nullable
- status (`draft`, `published`, `archived`)
- version_number
- previous_version_id nullable
- published_at nullable
- valid_from nullable
- valid_until nullable
- created_at
- updated_at

### diet_meals
- id
- tenant_id
- diet_id
- meal_name
- meal_order
- notes nullable

### diet_meal_items
- id
- tenant_id
- diet_meal_id
- food_item_id
- quantity_value
- quantity_unit
- household_measure_id nullable
- amount_description nullable
- preparation_notes nullable
- sort_order

### diet_substitutions
- id
- tenant_id
- diet_meal_item_id
- substitute_food_item_id
- quantity_value nullable
- quantity_unit nullable
- notes nullable
- sort_order

### diet_publications
- id
- tenant_id
- diet_id
- patient_id
- published_by_user_id
- published_at
- status

## 7. Documentos clínicos e exportações
### clinical_documents
- id
- tenant_id
- patient_id
- professional_id
- appointment_id nullable
- document_type (`exam_request`, `prescription`, `letter`, `other`)
- title
- status (`draft`, `finalized`, `published`)
- content_json
- version_number
- previous_version_id nullable
- created_at
- updated_at

### document_publications
- id
- tenant_id
- clinical_document_id
- patient_id
- published_by_user_id
- published_at

### exported_files
- id
- tenant_id
- related_entity_type
- related_entity_id
- export_type (`pdf`, `print_job`)
- file_key nullable
- status (`pending`, `processing`, `completed`, `failed`)
- requested_by_user_id
- created_at
- completed_at nullable
- failure_reason nullable

## 8. Notificações, suporte e auditoria
### notifications
- id
- tenant_id nullable
- recipient_user_id
- channel (`email`, `in_app`, `sms`, `whatsapp`, `push`)
- template_key
- payload_json
- status
- created_at
- sent_at nullable

### support_tickets
- id
- tenant_id nullable
- requester_user_id nullable
- category
- subject
- status
- priority
- assigned_backoffice_user_id nullable
- created_at
- updated_at

### support_ticket_events
- id
- support_ticket_id
- actor_user_id nullable
- actor_scope (`tenant`, `backoffice`, `system`)
- event_type
- payload_json
- created_at

### audit_logs
- id
- tenant_id nullable
- actor_user_id nullable
- actor_scope (`tenant`, `backoffice`, `system`)
- entity_type
- entity_id
- action
- reason nullable
- metadata_json
- ip_address nullable
- user_agent nullable
- created_at

## 9. Backoffice comercial e operacional
### backoffice_users
- id
- user_id
- is_active
- created_at

### sales_accounts
- id
- tenant_id nullable
- lead_name
- status
- owner_backoffice_user_id nullable
- metadata_json
- created_at
- updated_at

### tenant_operational_flags
- id
- tenant_id
- flag_key
- flag_value
- reason nullable
- created_by_user_id nullable
- created_at
