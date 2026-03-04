package models

type HealthBody struct {
	Status string `json:"status" example:"ok" doc:"Status of the server"`
}

type HealthOutput struct {
	Body HealthBody
}
