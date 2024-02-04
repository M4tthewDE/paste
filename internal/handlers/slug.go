package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/aws/aws-lambda-go/events"
	"github.com/m4tthewde/paste/internal"
)

func HandleSlug(ctx context.Context, req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	slug := req.PathParameters["slug"]
	log.Println(slug)

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Println("Empty bucket name")
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	c, err := internal.Download(ctx, bucketName, slug)
	if err != nil {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	formatter := formatters.Get("html")
	if formatter == nil {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	content := string(c)

	lexer := lexers.Analyse(content)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	style := styles.Get("github")
	if style == nil {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	var body bytes.Buffer
	err = formatter.Format(&body, style, iterator)
	if err != nil {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       body.String(),
	}, nil
}
