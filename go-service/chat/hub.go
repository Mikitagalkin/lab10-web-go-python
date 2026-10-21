package chat

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

// Message — сообщение чата.
type Message struct {
    Username string `json:"username"`
    Text     string `json:"text"`
}

// Client — клиент чата.
type Client struct {
    conn *websocket.Conn
    send chan Message
}

// Hub — центральный менеджер клиентов.
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan Message
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan Message),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Printf("Клиент подключился. Всего: %d", len(h.clients))
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                log.Printf("Клиент отключился. Всего: %d", len(h.clients))
            }
        case msg := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- msg:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}

// ServeWS — обработчик WebSocket.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{conn: conn, send: make(chan Message, 256)}
    hub.register <- client

    go client.writePump()
    client.readPump(hub)
}

func (c *Client) readPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()
    for {
        var msg Message
        if err := c.conn.ReadJSON(&msg); err != nil {
            break
        }
        hub.broadcast <- msg
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()
    for msg := range c.send {
        if err := c.conn.WriteJSON(msg); err != nil {
            return
        }
    }
}
