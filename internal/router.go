package internal

import (
	"bytes"
	"net/http"
	"os"
	"strconv"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		err := Index().Render(c.Request.Context(), c.Writer)
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

		slugLength := os.Getenv("SLUG_LENGTH")
		length, err := strconv.Atoi(slugLength)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			return
		}

		slug := Slug(length)

		err = Upload(c.Request.Context(), os.Getenv("BUCKET_NAME"), slug, content)
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

		contentBytes, err := Download(c.Request.Context(), os.Getenv("BUCKET_NAME"), slug)
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

		err = Paste(paste.String()).Render(c.Request.Context(), c.Writer)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
			return
		}
	})

	return r
}
