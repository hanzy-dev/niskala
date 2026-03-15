package domain

type IdempotencyStatus string

const (
	IdempotencyStatusProcessing IdempotencyStatus = "processing"
	IdempotencyStatusCompleted  IdempotencyStatus = "completed"
)

type IdempotencyRecord struct {
	UserID        string
	Key           string
	Status        IdempotencyStatus
	ResponseOrder *Order
}
