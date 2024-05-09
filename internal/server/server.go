package server

import (
	"github.com/Vaniog/Snaker/internal/game"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var hubRepo *HubRepository

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SetupRouter(r *gin.Engine) {
	hubRepo = NewHubRepository()

	r.GET("/find-hub/:mod", HandleFindHub)
	r.GET("/ws/play/:id", HandlePlay)
}

func HandleFindHub(c *gin.Context) {
	var opts game.Options
	switch c.Param("mod") {
	case "solo":
		opts = game.SoloOptions
	case "duo":
		opts = game.DuoOptions
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := hubRepo.FindHub(opts)
	resp := gin.H{
		"id": strconv.FormatInt(int64(id), 10),
	}

	c.JSON(http.StatusOK, resp)
}

func HandlePlay(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return
	}
	hub, err := hubRepo.GetHubById(HubID(id))
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	hub.RegisterClient(conn)
}
