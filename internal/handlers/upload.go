package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/m4tthewde/paste/internal"
)

func HandleUpload(ctx context.Context, req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	formData, err := url.ParseQuery(req.Body)
	if err != nil {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, err
	}

	content := formData.Get("content")
	if content == "" {
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "no content",
		}, errors.New("No content")
	}

	slugLength, err := getSlugLength()
	if err != nil {
		log.Println("Invalid slug length")
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	slug := internal.Slug(slugLength)

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Println("Empty bucket name")
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	err = internal.Upload(ctx, bucketName, slug, content)
	if err != nil {
		log.Println(err)
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       slug,
	}, nil
}

func getSlugLength() (int, error) {
	slugLength := os.Getenv("SLUG_LENGTH")
	length, err := strconv.Atoi(slugLength)
	return length, err

}
