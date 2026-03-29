package http

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Hub quản lý tất cả WebSocket clients và broadcast message
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mutex     sync.RWMutex
}

func NewHub() *Hub {
	h := &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte, 256),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for msg := range h.broadcast {
		h.mutex.RLock()
		for conn := range h.clients {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				conn.Close()
				// Upgrade to write lock để xóa client lỗi
				h.mutex.RUnlock()
				h.mutex.Lock()
				delete(h.clients, conn)
				h.mutex.Unlock()
				h.mutex.RLock()
			}
		}
		h.mutex.RUnlock()
	}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.mutex.Lock()
	delete(h.clients, conn)
	h.mutex.Unlock()
}

// Send gửi message đến tất cả clients đang kết nối
func (h *Hub) Send(msg []byte) {
	h.broadcast <- msg
}
