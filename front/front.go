package front

import (
	"net/http"
)

func init() {
	http.Handle("GET /", http.FileServer(http.Dir("front/static")))
}
