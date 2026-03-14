package schema

type GetHealthStatusRequest struct{}

type GetHealthStatusResponse struct {
	Body struct {
		Status string `json:"status" doc:"Status of the server" example:"ok"`
	}
}
