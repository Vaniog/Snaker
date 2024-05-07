package main

import (
	front "github.com/Vaniog/Snaker/front"
	backend "github.com/Vaniog/Snaker/internal/server"
	gin "github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	setupRouter(r)

	if err := r.Run(":8000"); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func setupRouter(r *gin.Engine) {
	front.SetupRouter(r)
	backend.SetupRouter(r)
}
