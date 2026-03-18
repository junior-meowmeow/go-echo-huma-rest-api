package entity

type Pet struct {
	Category  PetCategory
	ID        int64
	Name      string
	PhotoURLs []string
	Status    PetStatus
	Tags      []string
}

type PetCategory struct {
	ID   int64
	Name string
}

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusPending   PetStatus = "pending"
	PetStatusSold      PetStatus = "sold"
)
