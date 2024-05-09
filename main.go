package main

import (
	"flag"
	"fmt"
	front "github.com/Vaniog/Snaker/front"
	backend "github.com/Vaniog/Snaker/internal/server"
	gin "github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	setupRouter(r)

	port := flag.String("p", "8000", "port")
	flag.Parse()

	if err := r.Run(fmt.Sprintf(":%s", *port)); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func setupRouter(r *gin.Engine) {
	front.SetupRouter(r)
	backend.SetupRouter(r)
}
