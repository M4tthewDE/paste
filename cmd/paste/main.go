package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		err := internal.Index().Render(c.Request.Context(), c.Writer)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
		}
	})

	r.POST("/upload/paste", func(c *gin.Context) {

		content := c.PostForm("content")
		if content == "" {
			c.String(http.StatusBadRequest, "")
			return
		}

		slug := internal.Slug(config.SlugLength)

		err = internal.Upload(c.Request.Context(), config.BucketName, slug, content)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.String(http.StatusOK, slug)
	})

	r.GET("/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		if slug == "" {
			c.String(http.StatusBadRequest, "")
			return
		}

		contentBytes, err := internal.Download(c.Request.Context(), config.BucketName, slug)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
		}

		content := string(contentBytes)

		lexer := lexers.Analyse(content)
		if lexer == nil {
			lexer = lexers.Fallback
		}

		iterator, err := lexer.Tokenise(nil, content)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			return
		}

		style := styles.Get("github")
		if style == nil {
			c.String(http.StatusInternalServerError, "")
			return
		}

		formatter := html.New(html.Standalone(false), html.WithLineNumbers(true))
		if formatter == nil {
			c.String(http.StatusInternalServerError, "")
			return
		}

		var paste bytes.Buffer

		err = formatter.Format(&paste, style, iterator)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			return
		}

		err = internal.Paste(paste.String()).Render(c.Request.Context(), c.Writer)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			return
		}
	})

	log.Println("Listening on :" + config.Port)
	r.Run(":" + config.Port)
}
