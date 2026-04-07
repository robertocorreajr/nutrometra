package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClinicalDocument_Validate(t *testing.T) {
	tests := []struct {
		name    string
		doc     ClinicalDocument
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid document",
			doc:     ClinicalDocument{Title: "Solicitacao de exame", DocumentType: TypeExamRequest},
			wantErr: false,
		},
		{
			name:    "missing title fails",
			doc:     ClinicalDocument{DocumentType: TypePrescription},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name:    "invalid document_type fails",
			doc:     ClinicalDocument{Title: "Documento", DocumentType: "invalid_type"},
			wantErr: true,
			errMsg:  "invalid document_type",
		},
		{
			name:    "empty document_type fails",
			doc:     ClinicalDocument{Title: "Documento"},
			wantErr: true,
			errMsg:  "invalid document_type",
		},
		{
			name:    "exam_request type accepted",
			doc:     ClinicalDocument{Title: "Exame", DocumentType: TypeExamRequest},
			wantErr: false,
		},
		{
			name:    "prescription type accepted",
			doc:     ClinicalDocument{Title: "Prescricao", DocumentType: TypePrescription},
			wantErr: false,
		},
		{
			name:    "letter type accepted",
			doc:     ClinicalDocument{Title: "Carta", DocumentType: TypeLetter},
			wantErr: false,
		},
		{
			name:    "other type accepted",
			doc:     ClinicalDocument{Title: "Outro", DocumentType: TypeOther},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
