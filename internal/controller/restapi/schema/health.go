package schema

type GetHealthStatusRequest struct{}

type GetHealthStatusResponse struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Status of the server"`
	}
}
