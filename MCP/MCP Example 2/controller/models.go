package controller

type GreetInput struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type GreetOutput struct {
	Greeting string `json:"greeting"`
}
