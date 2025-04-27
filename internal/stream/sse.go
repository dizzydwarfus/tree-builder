package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dizzydwarfus/tree-builder/internal/shared"
	"github.com/dizzydwarfus/tree-builder/internal/treetraversal"
	"github.com/dizzydwarfus/tree-builder/types/trees"

	"github.com/google/uuid"
)

type Server struct {
	listenAddr string
	hub        *SSEHub
	dataCh     chan TreeParams
}

func (s *Server) dataProcessor() {
	log.Println("[Processor] Data processor started...")
	for params := range s.dataCh {
		log.Printf("[Processor] Processing data for session %s: %v\n", params.SessionId, params.Data)
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)

		err := treeInput(ctx, params, s.hub)
		if err != nil {
			log.Printf("[Processor] Error processing data for session %s: %v\n", params.SessionId, err)
		}
		log.Printf("[Processor] Successfully processed data for session %s\n", params.SessionId)
		cancel()
	}
	log.Println("[Processor] Data processor stopped...")
}

type SSEHub struct {
	mu      sync.Mutex
	streams map[string]chan trees.MultiChildTreeNode
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		streams: make(map[string]chan trees.MultiChildTreeNode),
	}
}

type TreeParams struct {
	SessionId string `json:"-"`
	Data      []int  `json:"data"`
}

func (hub *SSEHub) GetSessionChannel(sessionId string) chan trees.MultiChildTreeNode {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	ch, exists := hub.streams[sessionId]
	if !exists {
		ch = make(chan trees.MultiChildTreeNode, 50)
		hub.streams[sessionId] = ch
	}
	return ch
}

func (hub *SSEHub) CloseSessionChannel(sessionId string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if ch, ok := hub.streams[sessionId]; ok {
		close(ch)
		delete(hub.streams, sessionId)
	}
}

func (hub *SSEHub) Publish(sessionId string, data trees.MultiChildTreeNode) {
	ch := hub.GetSessionChannel(sessionId)
	ch <- data
}

func NewServer(listenaddr string) *Server {
	return &Server{
		listenAddr: listenaddr,
		hub:        NewSSEHub(),
		dataCh:     make(chan TreeParams),
	}
}

func (s *Server) Start() error {
	go s.dataProcessor()
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static/stream"))
	se := s.sendEvents()
	treePost := s.handleTreePost()

	mux.Handle("/", fs)
	mux.HandleFunc("/init", initHandler)
	mux.Handle("/tree", treePost)
	mux.Handle("/events", se)

	err := http.ListenAndServe(s.listenAddr, mux)
	log.Print(shared.Sred("[Server] HTTP server stopped.\n"))
	log.Print(shared.Sred("[Server] Closing data channel...\n"))
	close(s.dataCh)
	return err
}

func (s *Server) sendEvents() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Missing session cookie", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sessionId := cookie.Value
		// Prepare HTTP headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Flush the headers
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		flusher.Flush()

		ch := s.hub.GetSessionChannel(sessionId)
		// Keep pushing data as it arrives on the channel
		for {
			select {
			case data, ok := <-ch:
				if !ok {
					// The channel was closed. End this SSE connection.
					log.Printf(shared.Syellow("Channel closed for session %s"), sessionId)
					return
				}
				raw, err := json.Marshal(data)
				if err != nil {
					log.Printf("marshal: %v", err)
					continue
				}

				fmt.Fprintf(w, "data: %s\n\n", raw)
				flusher.Flush()

			case <-time.After(5 * time.Second):
				// keep-alive ping to avoid timeouts
				fmt.Fprintf(w, ": ping\n\n")
				flusher.Flush()

			case <-r.Context().Done():
				// Client disconnected
				log.Printf(shared.Sred("Client disconnected for session %s"), sessionId)
				// Optionally close & remove channel if no one else will use it
				s.hub.CloseSessionChannel(sessionId)
				return
			}
		}
	})
}

func (s *Server) handleTreePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Missing session cookie", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sessionId := cookie.Value
		var body map[string][]int

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		inputSlice, ok := body["data"]
		if !ok {
			http.Error(w, "Missing 'data' field", http.StatusBadRequest)
			return
		}
		log.Printf(shared.Sfaint("Received post request: %s\n"), inputSlice)

		select {
		case s.dataCh <- TreeParams{SessionId: sessionId, Data: inputSlice}:
			log.Printf(shared.Syellow("Data sent to processor for session %s\n"), sessionId)

		default:
			log.Printf(shared.Sred("Processor is busy, dropping data for session %s\n"), sessionId)
			http.Error(w, "Processor is busy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Accepted. Data sent to worker. SessionId: %s\n", sessionId)
	})
}

func ContentHandler(content string, contentType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write([]byte(content))
	}
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	// If no cookie, set a new one
	_, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		sessionId := uuid.NewString()
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: sessionId,
			Path:  "/",
		})
	}
	fmt.Fprintln(w, "OK, cookie set if missing")
}

func treeInput(ctx context.Context, data TreeParams, hub *SSEHub) error {
	root := trees.NewMultiChildTreeNode(1, "root", shared.Colors[0], 0)
	// check ctx err
	if ctx.Err() != nil {
		return fmt.Errorf("context error: %w", ctx.Err())
	}

	value := 2
	counter := &value
	treetraversal.TreeBuilder(root, data.Data, counter, 1)
	hub.Publish(data.SessionId, *root) // sends to the channel for that session

	return nil
}
