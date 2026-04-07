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
		audit.WithUserAgent("Mozilla/5.0"),
		audit.WithReason("signup"),
	)

	assert.Equal(t, tenantID, *entry.TenantID)
	assert.Equal(t, actorID, *entry.ActorUserID)
	assert.Equal(t, audit.ScopeTenant, entry.ActorScope)
	assert.Equal(t, "tenant", entry.EntityType)
	assert.Equal(t, entityID, *entry.EntityID)
	assert.Equal(t, "created", entry.Action)
	assert.NotEqual(t, uuid.Nil, entry.ID)
	assert.Equal(t, "192.168.1.1", entry.IPAddress)
	assert.Equal(t, "Mozilla/5.0", *entry.UserAgent)
	assert.Equal(t, "signup", *entry.Reason)
}

func TestWithIPAddress_StripsPort(t *testing.T) {
	entry := audit.NewEntry(audit.WithIPAddress("10.0.0.1:54321"))
	assert.Equal(t, "10.0.0.1", entry.IPAddress)
}

func TestWithIPAddress_InvalidIsIgnored(t *testing.T) {
	entry := audit.NewEntry(audit.WithIPAddress("not-an-ip"))
	assert.Equal(t, "", entry.IPAddress)
}

func TestNewService_NotNil(t *testing.T) {
	svc := audit.NewService()
	require.NotNil(t, svc)
}

func TestService_Write_NilExecutor(t *testing.T) {
	svc := audit.NewService()
	entry := audit.NewEntry(
		audit.WithActor(uuid.New(), audit.ScopeSystem),
		audit.WithEntity("test", uuid.New()),
		audit.WithAction("test"),
	)
	err := svc.Write(context.Background(), nil, entry)
	assert.Error(t, err)
}

func TestService_Write_EmptyActorScope(t *testing.T) {
	svc := audit.NewService()
	// Entry sem WithActor — ActorScope fica vazia
	entry := audit.NewEntry(
		audit.WithEntity("test", uuid.New()),
		audit.WithAction("test"),
	)
	err := svc.Write(context.Background(), nil, entry)
	assert.Error(t, err)
}
