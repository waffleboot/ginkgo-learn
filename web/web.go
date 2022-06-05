package web

import "github.com/google/uuid"

type MutableResponse struct {
	ServiceID   uuid.UUID `json:"service_id"`
	OperationID uuid.UUID `json:"operation_id"`
}

type Status string

var (
	StatusCreating = Status("creating")
	StatusRunning  = Status("running")
	StatusDeleting = Status("deleting")
	StatusDeleted  = Status("deleted")
)

type OperationResult struct {
	Status Status `json:"status"`
}
