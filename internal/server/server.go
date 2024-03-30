package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var hubRepo *Repository

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init() {
	hubRepo = NewHubRepository()

	http.HandleFunc("GET /find-hub/", ServeFindHub)
	http.HandleFunc("GET /ws/play/", ServePlay)
}

func ServeFindHub(w http.ResponseWriter, _ *http.Request) {
	id := hubRepo.FindHub()
	resp := map[string]string{
		"id": strconv.FormatInt(int64(id), 10),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func ServePlay(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/ws/play/"), 10, 64)
	if err != nil {
		return
	}
	hub, err := hubRepo.GetHubById(HubID(id))
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	player := newPlayer(hub, conn)
	hub.register <- player
}
