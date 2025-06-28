package realtime

import (
	"net/http"
	"pet-project/pkg/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	UserID int
	Conn   *websocket.Conn
	Send   chan model.Notification
}

type ClientManager struct {
	clients   map[int]*Client
	clientsMu sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[int]*Client),
	}
}

func (m *ClientManager) ServeWS(w http.ResponseWriter, r *http.Request, userID int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket", http.StatusBadRequest)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan model.Notification, 10),
	}

	m.clientsMu.Lock()
	if existingClient, exists := m.clients[userID]; exists {
		existingClient.Conn.Close()
	}
	m.clients[userID] = client
	m.clientsMu.Unlock()

	go m.writePump(client)
	go m.readPump(client)
}

func (m *ClientManager) writePump(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
		m.clientsMu.Lock()
		delete(m.clients, client.UserID)
		m.clientsMu.Unlock()
	}()

	for {
		select {
		case notif, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteJSON(notif); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (m *ClientManager) readPump(client *Client) {
	defer func() {
		client.Conn.Close()
		m.clientsMu.Lock()
		delete(m.clients, client.UserID)
		m.clientsMu.Unlock()
	}()

	client.Conn.SetReadLimit(512)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}
	}
}

func (m *ClientManager) Send(userID int, notif model.Notification) {
	m.clientsMu.RLock()
	client, ok := m.clients[userID]
	m.clientsMu.RUnlock()
	if ok {
		select {
		case client.Send <- notif:
		default:
		}
	}
}

func (m *ClientManager) Broadcast(notif model.Notification) {
	m.clientsMu.RLock()
	defer m.clientsMu.RUnlock()

	for _, client := range m.clients {
		select {
		case client.Send <- notif:
		default:
		}
	}
}

func (m *ClientManager) GetConnectedUsers() []int {
	m.clientsMu.RLock()
	defer m.clientsMu.RUnlock()

	users := make([]int, 0, len(m.clients))
	for userID := range m.clients {
		users = append(users, userID)
	}
	return users
}
