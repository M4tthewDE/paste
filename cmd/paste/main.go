package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/m4tthewde/paste/internal"
)

type Config struct {
	Data       string `json:"data"`
	SlugLength int    `json:"slugLength"`
}

var config Config

func parseConfig() error {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := parseConfig()
	if err != nil {
		log.Fatalln(err)
	}

	component := internal.Index()

	r := chi.NewRouter()
	r.Handle("/", templ.Handler(component))
	r.Post("/upload/paste", pasteHandler)
	r.Get("/{slug}", dataHandler)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func slug() string {
	b := make([]byte, config.SlugLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func pasteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slug := slug()
	fileName := config.Data + slug

	err = os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(slug))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	slugFile, err := os.Open(config.Data + slug)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	content, err := io.ReadAll(slugFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
