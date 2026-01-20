// internal/store/store.go
package store

//go:generate mockgen -destination=mock/mock.go -package=mock . MessageStore

import (
	"context"
	"errors"
	"time"
)

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

type MessageStore interface {
	FindRecipient(ctx context.Context, username string) (userID string, err error)
	ListMessages(ctx context.Context, userID string) ([]Message, error)
	GetMessage(ctx context.Context, id int64) (*Message, error)
	SaveMessage(ctx context.Context, userID string, msg Message) error
	// RegisterUser регистрирует нового пользователя
	RegisterUser(ctx context.Context, userID, username string) error
}

// Message представляет сообщение в системе
type Message struct {
	ID      int64     `json:"id"`
	Sender  string    `json:"sender"`
	Time    time.Time `json:"time"`
	Payload string    `json:"payload"`
}
