package event_manager

import (
	"slices"
	"sync"
)

type EventData struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

// SubscriberInfo stores information about a subscriber
type SubscriberInfo struct {
	Channel  chan EventData
	UserID   int
	TenantID string
}

// EventManager manages subscribers and sends events
// It is a singleton
// Once a client is subscribed to a channel, it will receive all events for that channel
type EventManager struct {
	subscribers map[string][]SubscriberInfo // channel -> list of subscribers
	mutex       sync.RWMutex
}

var (
	instance *EventManager
	once     sync.Once
)

// NewEventManager creates a new event manager instance (singleton)
func NewEventManager() *EventManager {
	once.Do(func() {
		instance = &EventManager{
			subscribers: make(map[string][]SubscriberInfo),
		}
	})
	return instance
}

// Subscribe subscribes a client to a channel
func (em *EventManager) Subscribe(channel string, userID int, tenantID string) chan EventData {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	ch := make(chan EventData, 1000)
	em.subscribers[channel] = append(em.subscribers[channel], SubscriberInfo{
		Channel:  ch,
		UserID:   userID,
		TenantID: tenantID,
	})
	return ch
}

// Unsubscribe unsubscribes a client from a channel
func (em *EventManager) Unsubscribe(channel string, ch chan EventData) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if subs, found := em.subscribers[channel]; found {
		for i, subscriber := range subs {
			if subscriber.Channel == ch {
				em.subscribers[channel] = slices.Delete(subs, i, i+1)
				close(ch)
				break
			}
		}
		if len(em.subscribers[channel]) == 0 {
			delete(em.subscribers, channel)
		}
	}
}

// SendEvent sends an event to all subscribers of a channel
func (em *EventManager) SendEvent(channel, event string, data any) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	if subs, found := em.subscribers[channel]; found {
		eventData := EventData{
			Event: event,
			Data:  data,
		}
		for _, sub := range subs {
			select {
			case sub.Channel <- eventData:
			default:
				// Skip if the channel is busy
			}
		}
	}
}
