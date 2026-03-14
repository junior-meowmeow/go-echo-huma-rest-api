package schema

type Pet struct {
	ID        int64       `json:"id" doc:"Pet ID"`
	Name      string      `json:"name" doc:"Pet name"`
	Category  PetCategory `json:"category" doc:"Pet category"`
	PhotoURLs []string    `json:"photoUrls" doc:"Pet Photo URLs"`
	Status    string      `json:"status" enum:"available,pending,sold" doc:"Pet status"`
	Tags      []string    `json:"tags" doc:"Pet tags"`
}

type PetCategory struct {
	ID   int64  `json:"id" doc:"Category ID"`
	Name string `json:"name" doc:"Category name"`
}

type GetAvailablePetsRequest struct{}

type GetAvailablePetsResponse struct {
	Body struct {
		Data []Pet `json:"data" doc:"List of available pets"`
	}
}

type GetPetByIDRequest struct {
	ID int64 `path:"id" required:"true" doc:"Pet ID"`
}

type GetPetByIDResponse struct {
	Body Pet
}
