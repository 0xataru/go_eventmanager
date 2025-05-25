package event_manager

import (
	"testing"
	"time"
)

func TestEventManager_Subscribe(t *testing.T) {
	em := NewEventManager()
	channel := "test-channel"
	userID := 1
	tenantID := "test-tenant"

	// Test subscribing to a channel
	ch := em.Subscribe(channel, userID, tenantID)

	// Verify that the channel was added to subscribers
	em.mutex.RLock()
	if len(em.subscribers[channel]) != 1 {
		t.Errorf("Expected 1 subscriber, got %d", len(em.subscribers[channel]))
	}
	em.mutex.RUnlock()

	// Test sending and receiving an event
	go func() {
		em.SendEvent(channel, "test-event", "test-data")
	}()

	select {
	case event := <-ch:
		if event.Event != "test-event" {
			t.Errorf("Expected event 'test-event', got '%s'", event.Event)
		}
		if event.Data != "test-data" {
			t.Errorf("Expected data 'test-data', got '%v'", event.Data)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for event")
	}

	// Clean up
	em.Unsubscribe(channel, ch)
}

func TestEventManager_Unsubscribe(t *testing.T) {
	em := NewEventManager()
	channel := "test-channel"

	// Subscribe to a channel
	ch := em.Subscribe(channel, 1, "test-tenant")

	// Unsubscribe from the channel
	em.Unsubscribe(channel, ch)

	// Verify that the channel was removed from subscribers
	em.mutex.RLock()
	if _, exists := em.subscribers[channel]; exists {
		t.Error("Channel should be removed from subscribers")
	}
	em.mutex.RUnlock()

	// Verify that the channel is closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("Channel should be closed")
		}
	default:
		t.Error("Channel should be closed")
	}
}

func TestEventManager_SendEvent(t *testing.T) {
	em := NewEventManager()
	channel := "test-channel"
	eventCount := 3

	// Create multiple subscribers
	channels := make([]chan EventData, eventCount)
	for i := 0; i < eventCount; i++ {
		channels[i] = em.Subscribe(channel, 1, "test-tenant")
	}

	// Send an event
	em.SendEvent(channel, "test-event", "test-data")

	// Verify that all subscribers received the event
	for i, ch := range channels {
		select {
		case event := <-ch:
			if event.Event != "test-event" {
				t.Errorf("Subscriber %d: Expected event 'test-event', got '%s'", i, event.Event)
			}
			if event.Data != "test-data" {
				t.Errorf("Subscriber %d: Expected data 'test-data', got '%v'", i, event.Data)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("Subscriber %d: Timeout waiting for event", i)
		}
	}

	// Clean up
	for _, ch := range channels {
		em.Unsubscribe(channel, ch)
	}
}
