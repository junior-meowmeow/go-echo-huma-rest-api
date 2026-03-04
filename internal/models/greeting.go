package models

type GreetingInput struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

type GreetingOutputBody struct {
	Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
}

type GreetingOutput struct {
	Body GreetingOutputBody
}
