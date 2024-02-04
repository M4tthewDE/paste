package handlers

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/m4tthewde/paste/internal"
)

func HandleIndex(ctx context.Context) (*events.APIGatewayProxyResponse, error) {
	component := internal.Index()

	var body bytes.Buffer
	err := component.Render(ctx, &body)
	if err != nil {
		return nil, err
	}

	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body.String(),
	}

	return &res, nil
}
