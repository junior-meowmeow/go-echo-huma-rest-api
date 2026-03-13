package schema

type GreetingRequest struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

type GreetingResponse struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}
