package web

import "github.com/google/uuid"

type MutableResponse struct {
	ServiceID   uuid.UUID
	OperationID uuid.UUID
}

type Status string

var (
	StatusCreating = Status("creating")
	StatusRunning  = Status("running")
	StatusDeleting = Status("deleting")
	StatusDeleted  = Status("deleted")
)

type OperationResult struct {
	Status Status
}
