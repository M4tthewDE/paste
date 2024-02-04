package main

import (
	"log"

	"github.com/m4tthewde/paste/internal"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := internal.Router()
	log.Println("Listening on :8080")
	r.Run(":8080")
}
