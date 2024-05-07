package front

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static/*
var static embed.FS

func SetupRouter(r *gin.Engine) {
	err := godotenv.Load()
	useEmbed := true
	if err == nil {
		if envDebug, has := os.LookupEnv("DEBUG"); has {
			if envDebug == "True" {
				useEmbed = false
			}
		}
	}

	var staticFS http.FileSystem
	if useEmbed {
		serverRoot, err := fs.Sub(static, "static")
		if err != nil {
			log.Fatal(err)
		}
		staticFS = http.FS(serverRoot)
	} else {
		staticFS = http.Dir("front/static")
	}

	r.NoRoute(gin.WrapH(http.FileServer(staticFS)))
}
