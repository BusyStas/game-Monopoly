package main

import (
    "embed"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

//go:embed static/*
var staticFS embed.FS

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Hub struct {
    mu      sync.Mutex
    clients map[*websocket.Conn]bool
    broadcast chan []byte
}

func NewHub() *Hub {
    return &Hub{
        clients:   make(map[*websocket.Conn]bool),
        broadcast: make(chan []byte, 16),
    }
}

func (h *Hub) Run() {
    for msg := range h.broadcast {
        h.mu.Lock()
        for c := range h.clients {
            _ = c.WriteMessage(websocket.TextMessage, msg)
        }
        h.mu.Unlock()
    }
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("ws upgrade:", err)
        return
    }
    h.mu.Lock()
    h.clients[c] = true
    h.mu.Unlock()

    // read loop
    go func(cc *websocket.Conn) {
        defer func() {
            h.mu.Lock()
            delete(h.clients, cc)
            h.mu.Unlock()
            cc.Close()
        }()
        for {
            mt, p, err := cc.ReadMessage()
            if err != nil {
                return
            }
            if mt == websocket.TextMessage {
                // broadcast to all
                h.broadcast <- p
            }
        }
    }(c)
}

func main() {
    hub := NewHub()
    go hub.Run()

    // HTTP handlers
    fs := http.FS(staticFS)
    staticHandler := http.FileServer(fs)

    http.Handle("/", http.StripPrefix("/", staticHandler))
    http.HandleFunc("/ws", hub.ServeWS)

    addr := ":8080"
    log.Printf("Starting server on %s — open http://localhost:8080/ in your browser", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatal(err)
    }
}
