package main

import (
	"github.com/bluegopher/alice-skill/internal/store"
)

type app struct {
	store store.MessageStore
}

// NewApp создает новый экземпляр приложения.
func NewApp(store store.MessageStore) *app {
	return &app{store: store}
}