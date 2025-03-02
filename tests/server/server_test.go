package server_test

import (
	"pubsub/internal/server"
	"testing"
	"time"
)

func TestPublishSubscribe(t *testing.T) {
	s := server.NewServer()
	defer s.Close()

	sub := s.Subscribe("food")
	s.PublishMessage("food", "Hello, Food!")

	select {
	case msg := <-sub:
		if msg != "Hello, Food!" {
			t.Errorf("Expected 'Hello, Food!', but got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestClosedServerDoesNotPublish(t *testing.T) {
	s := server.NewServer()
	sub := s.Subscribe("sports")
	s.Close()

	// Try publishing after server is closed
	s.PublishMessage("sports", "No one should get this")

	_, ok := <-sub
	if ok {
		t.Error("Expected no message, but received a message from a closed channel")
	}
}

func TestMultipleSubscribers(t *testing.T) {
	s := server.NewServer()
	defer s.Close()

	sub1 := s.Subscribe("news")
	sub2 := s.Subscribe("news")

	s.PublishMessage("news", "Breaking News!")

	select {
	case msg := <-sub1:
		if msg != "Breaking News!" {
			t.Errorf("Expected 'Breaking News!', got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for message on sub1")
	}

	select {
	case msg := <-sub2:
		if msg != "Breaking News!" {
			t.Errorf("Expected 'Breaking News!', got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for message on sub2")
	}
}

func TestMessageOrder(t *testing.T) {
	s := server.NewServer()
	defer s.Close()

	sub := s.Subscribe("chat")

	messages := []string{"Hello", "How are you?", "Goodbye"}
	for _, msg := range messages {
		s.PublishMessage("chat", msg)
	}

	for _, expected := range messages {
		select {
		case msg := <-sub:
			if msg != expected {
				t.Errorf("Expected '%s', got '%s'", expected, msg)
			}
		case <-time.After(time.Second):
			t.Errorf("Timeout waiting for message '%s'", expected)
		}
	}
}
