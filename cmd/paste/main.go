package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/go-chi/chi/v5"
	"github.com/m4tthewde/paste/internal"
)

var config *internal.Config

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c, err := internal.ParseConfig()
	if err != nil {
		log.Fatalln(err)
	}

	config = c

	r := chi.NewRouter()
	r.Handle("/", templ.Handler(internal.Index()))
	r.Post("/upload/paste", uploadHandler)
	r.Get("/{slug}", slugHandler)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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

	slug := internal.Slug(config.SlugLength)
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

func slugHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	slugFile, err := os.Open(config.Data + slug)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c, err := io.ReadAll(slugFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Has("raw") {
		_, err = w.Write(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	content := string(c)

	lexer := lexers.Analyse(content)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	style := styles.Get("github")
	if style == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	formatter := formatters.Get("html")
	if formatter == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = formatter.Format(w, style, iterator)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
