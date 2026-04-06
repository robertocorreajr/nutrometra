package audit_test

import (
	"context"
	"testing"

	"nutrometra/api/internal/platform/audit"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEntry(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	entityID := uuid.New()

	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("tenant", entityID),
		audit.WithAction("created"),
		audit.WithIPAddress("192.168.1.1"),
	)

	assert.Equal(t, tenantID, *entry.TenantID)
	assert.Equal(t, actorID, *entry.ActorUserID)
	assert.Equal(t, audit.ScopeTenant, entry.ActorScope)
	assert.Equal(t, "tenant", entry.EntityType)
	assert.Equal(t, entityID, *entry.EntityID)
	assert.Equal(t, "created", entry.Action)
	assert.NotEqual(t, uuid.Nil, entry.ID)
}

func TestService_Write_RequiresDB(t *testing.T) {
	// Smoke test: New retorna sem pânico
	svc := audit.NewService(nil)
	require.NotNil(t, svc)

	// Write com db nil deve retornar erro, não panicar
	entry := audit.NewEntry(
		audit.WithActor(uuid.New(), audit.ScopeSystem),
		audit.WithEntity("test", uuid.New()),
		audit.WithAction("test"),
	)
	err := svc.Write(context.Background(), entry)
	assert.Error(t, err)
}
