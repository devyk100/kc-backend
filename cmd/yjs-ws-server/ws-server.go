package ws_server

import (
	"encoding/binary"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Yjs message types
const (
	messageSync      = 0
	messageAwareness = 1
)

// Document represents a shared Yjs document
type Document struct {
	Name        string
	Connections map[*websocket.Conn]bool
	Awareness   map[uint32]interface{}
	Content     []byte
	Mutex       sync.Mutex
}

// Server holds the Yjs server state
type Server struct {
	Documents map[string]*Document
	Upgrader  websocket.Upgrader
	Mutex     sync.Mutex
}

func Start() {
	// Initialize the server
	server := &Server{
		Documents: make(map[string]*Document),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	// Start the WebSocket server
	http.HandleFunc("/", server.handleWebSocket)
	log.Println("Yjs WebSocket server started on :1234")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

// handleWebSocket handles incoming WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Extract the document name from the URL path
	docName := r.URL.Path[1:]
	if docName == "" {
		log.Println("No document name provided")
		return
	}

	// Get or create the document
	doc := s.getOrCreateDocument(docName)
	doc.Mutex.Lock()
	doc.Connections[conn] = true
	doc.Mutex.Unlock()

	// Handle incoming messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Only process binary messages
		if messageType != websocket.BinaryMessage {
			log.Println("Received non-binary message, ignoring")
			continue
		}

		// Process the Yjs message
		err = s.handleYjsMessage(conn, doc, message)
		if err != nil {
			log.Println("Yjs message error:", err)
			break
		}
	}

	// Clean up when the connection closes
	doc.Mutex.Lock()
	delete(doc.Connections, conn)
	doc.Mutex.Unlock()
}

// getOrCreateDocument retrieves or creates a document
func (s *Server) getOrCreateDocument(name string) *Document {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if doc, exists := s.Documents[name]; exists {
		return doc
	}

	doc := &Document{
		Name:        name,
		Connections: make(map[*websocket.Conn]bool),
		Awareness:   make(map[uint32]interface{}),
	}
	s.Documents[name] = doc
	return doc
}

// handleYjsMessage processes incoming Yjs messages
func (s *Server) handleYjsMessage(conn *websocket.Conn, doc *Document, message []byte) error {
	if len(message) < 1 {
		return errors.New("empty message")
	}

	// The first byte is the message type
	messageType := message[0]
	switch messageType {
	case messageSync:
		return s.handleSyncMessage(conn, doc, message[1:])
	case messageAwareness:
		return s.handleAwarenessMessage(conn, doc, message[1:])
	default:
		return errors.New("unknown message type")
	}
}

// handleSyncMessage processes Yjs sync protocol messages
func (s *Server) handleSyncMessage(conn *websocket.Conn, doc *Document, message []byte) error {
	doc.Mutex.Lock()
	defer doc.Mutex.Unlock()

	// Broadcast the sync message to all connected clients
	for client := range doc.Connections {
		if client != conn {
			err := client.WriteMessage(websocket.BinaryMessage, append([]byte{messageSync}, message...))
			if err != nil {
				log.Println("Write error:", err)
			}
		}
	}

	// Update the document content
	doc.Content = message
	return nil
}

// handleAwarenessMessage processes Yjs awareness protocol messages
func (s *Server) handleAwarenessMessage(conn *websocket.Conn, doc *Document, message []byte) error {
	// Decode the awareness message (simplified for demonstration)
	if len(message) < 4 {
		return errors.New("invalid awareness message")
	}

	// The first 4 bytes represent the client ID
	clientID := binary.BigEndian.Uint32(message[:4])
	awarenessUpdate := message[4:]

	doc.Mutex.Lock()
	defer doc.Mutex.Unlock()

	// Update the awareness state
	doc.Awareness[clientID] = awarenessUpdate

	// Broadcast the awareness update to all connected clients
	for client := range doc.Connections {
		if client != conn {
			err := client.WriteMessage(websocket.BinaryMessage, append([]byte{messageAwareness}, message...))
			if err != nil {
				log.Println("Write error:", err)
			}
		}
	}

	return nil
}
