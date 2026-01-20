package main

import (
	"github.com/bluegopher/alice-skill/internal/store"
)

type app struct {
	store store.MessageStore
}

// newApp создает новый экземпляр приложения
func newApp(store store.MessageStore) *app {
	return &app{store: store}
}
