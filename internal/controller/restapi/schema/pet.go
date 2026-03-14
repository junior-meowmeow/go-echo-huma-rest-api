package schema

type Pet struct {
	ID        int64       `json:"id" doc:"Pet ID"`
	Name      string      `json:"name" doc:"Pet name"`
	Category  PetCategory `json:"category" doc:"Pet category"`
	PhotoURLs []string    `json:"photoUrls" doc:"Pet Photo URLs"`
	Status    string      `json:"status" doc:"Pet status (available, pending, sold)"`
	Tags      []string    `json:"tags" doc:"Pet tags"`
}

type PetCategory struct {
	ID   int64  `json:"id" doc:"Category ID"`
	Name string `json:"name" doc:"Category name"`
}

type GetAvailablePetsRequest struct{}

type GetAvailablePetsResponse struct {
	Body struct {
		Data []Pet `json:"data"`
	}
}

type GetPetByIDRequest struct {
	ID int64 `path:"id" required:"true" doc:"Pet ID"`
}

type GetPetByIDResponse struct {
	Body Pet
}
