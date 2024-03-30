package main

import (
	_ "github.com/Vaniog/Snaker/front"
	_ "github.com/Vaniog/Snaker/internal/server"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
