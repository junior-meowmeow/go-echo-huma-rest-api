package schema

type GreetingRequest struct {
	Name string `path:"name" maxLength:"30" doc:"Name to greet" example:"world"`
}

type GreetingResponse struct {
	Body struct {
		Message string `json:"message" doc:"Greeting message" example:"Hello, world!"`
	}
}
