package domain

import "github.com/google/uuid"

// Role represents an RBAC role scoped to an application area.
type Role struct {
	ID               uuid.UUID
	Code             string
	ApplicationScope string
	Name             string
	Description      string
}

// Permission represents a granular action that can be assigned to roles.
type Permission struct {
	ID               uuid.UUID
	Code             string
	ApplicationScope string
	Description      string
}
