package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfessional_Validate(t *testing.T) {
	tests := []struct {
		name    string
		p       Professional
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid CRN with state",
			p: Professional{
				FullName:           "Ana Silva",
				RegistrationType:   RegistrationCRN,
				RegistrationNumber: "12345",
				RegistrationState:  "SP",
			},
			wantErr: false,
		},
		{
			name: "CRN without state fails",
			p: Professional{
				FullName:           "Ana Silva",
				RegistrationType:   RegistrationCRN,
				RegistrationNumber: "12345",
			},
			wantErr: true,
			errMsg:  "registration_state is required",
		},
		{
			name: "CRM without state fails",
			p: Professional{
				FullName:           "Dr. João",
				RegistrationType:   RegistrationCRM,
				RegistrationNumber: "67890",
			},
			wantErr: true,
			errMsg:  "registration_state is required",
		},
		{
			name: "other type without state passes",
			p: Professional{
				FullName:           "Carlos",
				RegistrationType:   RegistrationOther,
				RegistrationNumber: "OTHER-001",
			},
			wantErr: false,
		},
		{
			name: "missing full_name fails",
			p: Professional{
				RegistrationType:   RegistrationCRN,
				RegistrationNumber: "12345",
				RegistrationState:  "SP",
			},
			wantErr: true,
			errMsg:  "full_name is required",
		},
		{
			name: "missing registration_number fails",
			p: Professional{
				FullName:         "Ana Silva",
				RegistrationType: RegistrationCRN,
			},
			wantErr: true,
			errMsg:  "registration_number is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.p.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
