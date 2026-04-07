CREATE TABLE clinical_documents (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id),
    patient_id          UUID NOT NULL REFERENCES patients(id),
    professional_id     UUID NOT NULL REFERENCES professionals(id),
    appointment_id      UUID REFERENCES appointments(id),
    document_type       TEXT NOT NULL
                        CHECK (document_type IN ('exam_request', 'prescription', 'letter', 'other')),
    title               TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'finalized', 'published')),
    content_json        JSONB NOT NULL DEFAULT '{}',
    version_number      INT NOT NULL DEFAULT 1,
    previous_version_id UUID REFERENCES clinical_documents(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_clin_docs_patient ON clinical_documents(tenant_id, patient_id);

CREATE TABLE document_publications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants(id),
    clinical_document_id    UUID NOT NULL REFERENCES clinical_documents(id),
    patient_id              UUID NOT NULL REFERENCES patients(id),
    published_by_user_id    UUID NOT NULL REFERENCES users(id),
    published_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_doc_pub_patient ON document_publications(tenant_id, patient_id);
