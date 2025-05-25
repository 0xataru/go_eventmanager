# Go Event Manager

[![Go Report Card](https://goreportcard.com/badge/github.com/0xataru/go_eventmanager)](https://goreportcard.com/report/github.com/0xataru/go_eventmanager)
[![GoDoc](https://godoc.org/github.com/0xataru/go_eventmanager?status.svg)](https://godoc.org/github.com/0xataru/go_eventmanager)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/0xataru/go_eventmanager/actions)

A lightweight, thread-safe event management system for Go applications that implements the publish-subscribe pattern.

## Features

- Thread-safe event management
- Singleton pattern implementation
- Support for multiple channels
- User and tenant-based subscription system
- Buffered channels for event delivery
- Automatic cleanup of empty channels

## Installation

```bash
go get github.com/0xataru/go_eventmanager
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/0xataru/go_eventmanager"
)

func main() {
    // Get the event manager instance
    em := event_manager.NewEventManager()

    // Subscribe to a channel
    ch := em.Subscribe("my-channel", 1, "tenant-1")

    // Start listening for events
    go func() {
        for event := range ch {
            fmt.Printf("Received event: %s with data: %v\n", event.Event, event.Data)
        }
    }()

    // Send an event
    em.SendEvent("my-channel", "user-login", map[string]string{
        "username": "john_doe",
    })

    // Unsubscribe when done
    em.Unsubscribe("my-channel", ch)
}
```

### API Reference

#### EventManager

The main struct that manages all event subscriptions and deliveries.

```go
type EventManager struct {
    subscribers map[string][]SubscriberInfo
    mutex       sync.RWMutex
}
```

#### Methods

- `NewEventManager() *EventManager` - Creates a new event manager instance (singleton)
- `Subscribe(channel string, userID int, tenantID string) chan EventData` - Subscribes to a channel
- `Unsubscribe(channel string, ch chan EventData)` - Unsubscribes from a channel
- `SendEvent(channel, event string, data any)` - Sends an event to all subscribers of a channel

#### EventData

```go
type EventData struct {
    Event string `json:"event"`
    Data  any    `json:"data"`
}
```

## Best Practices

1. Always unsubscribe when you're done with a channel to prevent memory leaks
2. Use buffered channels (default size is 1000) to handle high-frequency events
3. Consider implementing error handling for your event processing logic
4. Use meaningful channel names that reflect your application's domain

## Thread Safety

The EventManager is fully thread-safe and can be used in concurrent applications. It uses a read-write mutex to protect the subscribers map and ensure safe concurrent access.

## License

This project is dual-licensed under both the MIT License and Apache License 2.0.