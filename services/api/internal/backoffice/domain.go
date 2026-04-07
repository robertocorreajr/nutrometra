package backoffice

import (
	"time"

	"github.com/google/uuid"
)

// BackofficeUser represents an admin/support user with backoffice access.
type BackofficeUser struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	IsActive  bool
	CreatedAt time.Time
}
