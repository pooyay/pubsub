package pubsub

import "sync"

type Server struct {
	mu     sync.Mutex
	subs   map[string][]chan string
	quit   chan struct{}
	closed bool
}

func NewServer() *Server {
	return &Server{
		subs: make(map[string][]chan string),
		quit: make(chan struct{}),
	}
}

func (s *Server) PublishMessage(topic string, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}

	for _, ch := range s.subs[topic] {
		select {
		case ch <- message:
		default:
		}

	}
}

func (s *Server) Subscribe(topic string) <-chan string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		ch := make(chan string)
		close(ch) // Return a closed channel instead of nil to prevent runtime panics
		return ch
	}

	ch := make(chan string, 10)
	s.subs[topic] = append(s.subs[topic], ch)
	return ch
}

func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}

	s.closed = true
	close(s.quit)

	for _, ch := range s.subs {
		for _, subscriber := range ch {
			close(subscriber)
		}
	}
}
