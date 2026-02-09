package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"projectwebcurhat/contract"
	"projectwebcurhat/database"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketController struct {
	service *contract.Service
}

func (w *WebSocketController) GetPrefix() string {
	return "/ws"
}

func (w *WebSocketController) InitService(service *contract.Service) {
	w.service = service
}

func (w *WebSocketController) InitRoute(app *gin.RouterGroup) {
	app.GET("", w.HandleConnection)
}

func (w *WebSocketController) HandleConnection(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	username := ctx.Query("username")
	if username == "" {
		username = "Anonymous"
	}

	clientID := uuid.New().String()
	client := database.NewClient(clientID, conn, username)

	log.Printf("New client connected: %s (username: %s)", clientID, username)

	go w.writePump(client)
	go w.readPump(client)
}

func (w *WebSocketController) readPump(client *database.Client) {
	defer func() {
		w.service.Signaling.DisconnectClient(client)
		client.Conn.Close()
		log.Printf("Client %s disconnected", client.ID)
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		if err := w.service.Signaling.HandleMessage(client, message); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}

func (w *WebSocketController) writePump(client *database.Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		message, ok := <-client.Send
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		writer, err := client.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		writer.Write(message)

		n := len(client.Send)
		for i := 0; i < n; i++ {
			writer.Write([]byte{'\n'})
			writer.Write(<-client.Send)
		}

		if err := writer.Close(); err != nil {
			return
		}
	}
}
