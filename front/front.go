package front

import (
	"embed"
	"github.com/joho/godotenv"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static/*
var static embed.FS

func init() {
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

	http.Handle("GET /", http.FileServer(staticFS))
}
