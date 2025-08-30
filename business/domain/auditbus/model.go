package auditbus

import (
	"encoding/json"
	"service/business/types/domain"
	"service/business/types/name"
	"time"

	"github.com/google/uuid"
)

// Audit represents information about an individual audit record.
type Audit struct {
	ID        uuid.UUID
	ObjID     uuid.UUID
	ObjDomain domain.Domain
	ObjName   name.Name
	ActorID   uuid.UUID
	Action    string
	Data      json.RawMessage
	Message   string
	Timestamp time.Time
}

// NewAudit represents the information needed to create a new audit record.
type NewAudit struct {
	ObjID     uuid.UUID
	ObjDomain domain.Domain
	ObjName   name.Name
	ActorID   uuid.UUID
	Action    string
	Data      any
	Message   string
}
