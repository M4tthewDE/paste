package main

import (
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-chi/chi/v5"
	"github.com/m4tthewde/paste/internal/handlers"
)

// only for local development
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()
	r.Get("/", indexHandler)
	r.Post("/upload/paste", uploadHandler)
	r.Get("/{slug}", slugHandler)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := handlers.HandleIndex(r.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != 200 {
		w.WriteHeader(resp.StatusCode)
		return
	}

	_, err = w.Write([]byte(resp.Body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := events.APIGatewayV2HTTPRequest{
		Body: string(body),
	}

	resp, err := handlers.HandleUpload(r.Context(), req)
	if resp.StatusCode != 200 {
		w.WriteHeader(resp.StatusCode)
		return
	}

	_, err = w.Write([]byte(resp.Body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func slugHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		log.Println("TEST")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"slug": slug},
	}

	resp, err := handlers.HandleSlug(r.Context(), req)
	if resp.StatusCode != 200 {
		w.WriteHeader(resp.StatusCode)
		return
	}

	_, err = w.Write([]byte(resp.Body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
