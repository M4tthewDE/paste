package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/m4tthewde/paste/internal/handlers"
)

func main() {
	lambda.Start(handlers.HandleRequest)
}
