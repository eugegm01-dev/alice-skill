package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bluegopher/alice-skill/internal/models"
	"github.com/bluegopher/alice-skill/internal/store"
)

// inMemoryStore — простейшая реализация хранилища для примеров.
type inMemoryStore struct {
	users    map[string]string // userID -> username
	messages map[string][]store.Message
}

func (s *inMemoryStore) FindRecipient(ctx context.Context, username string) (string, error) {
	for uid, uname := range s.users {
		if uname == username {
			return uid, nil
		}
	}
	return "", errors.New("user not found")
}

func (s *inMemoryStore) ListMessages(ctx context.Context, userID string) ([]store.Message, error) {
	msgs, ok := s.messages[userID]
	if !ok {
		return []store.Message{}, nil
	}
	return msgs, nil
}

func (s *inMemoryStore) GetMessage(ctx context.Context, id int64) (*store.Message, error) {
	for _, msgs := range s.messages {
		for _, msg := range msgs {
			if msg.ID == id {
				return &msg, nil
			}
		}
	}
	return nil, errors.New("message not found")
}

func (s *inMemoryStore) SaveMessage(ctx context.Context, userID string, msg store.Message) error {
	msg.ID = int64(len(s.messages[userID]) + 1)
	s.messages[userID] = append(s.messages[userID], msg)
	return nil
}

func (s *inMemoryStore) RegisterUser(ctx context.Context, userID, username string) error {
	if _, exists := s.users[userID]; exists {
		return nil
	}
	for _, uname := range s.users {
		if uname == username {
			return store.ErrConflict
		}
	}
	s.users[userID] = username
	return nil
}

// sendPost — вспомогательная функция для отправки POST-запроса.
func sendPost(url string, req models.Request) []byte {
	body, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)
	return result
}

// Example_registerUser демонстрирует регистрацию пользователя.
func Example_registerUser() {
	memStore := &inMemoryStore{
		users:    make(map[string]string),
		messages: make(map[string][]store.Message),
	}
	appInstance := NewApp(memStore)

	ts := httptest.NewServer(http.HandlerFunc(appInstance.webhook))
	defer ts.Close()

	req := models.Request{
		Request: models.SimpleUtterance{
			Type:    models.TypeSimpleUtterance,
			Command: "Зарегистрируй Alice",
		},
		Session: models.Session{
			New:  false,
			User: models.RequestUser{UserID: "user1"},
		},
		Version: "1.0",
	}
	_ = sendPost(ts.URL, req)
	// В ответе ожидается "Вы успешно зарегистрированы под именем Alice"
}

// Example_sendMessage демонстрирует отправку сообщения.
func Example_sendMessage() {
	memStore := &inMemoryStore{
		users:    make(map[string]string),
		messages: make(map[string][]store.Message),
	}
	_ = memStore.RegisterUser(context.Background(), "alice_id", "Alice")

	appInstance := NewApp(memStore)
	ts := httptest.NewServer(http.HandlerFunc(appInstance.webhook))
	defer ts.Close()

	req := models.Request{
		Request: models.SimpleUtterance{
			Type:    models.TypeSimpleUtterance,
			Command: "Отправь Alice Привет, как дела?",
		},
		Session: models.Session{
			New:  false,
			User: models.RequestUser{UserID: "user2"},
		},
		Version: "1.0",
	}
	_ = sendPost(ts.URL, req)
	// Ожидается ответ "Сообщение успешно отправлено"
}

// Example_readMessage демонстрирует чтение сообщения.
func Example_readMessage() {
	memStore := &inMemoryStore{
		users:    make(map[string]string),
		messages: make(map[string][]store.Message),
	}
	_ = memStore.RegisterUser(context.Background(), "alice_id", "Alice")
	_ = memStore.RegisterUser(context.Background(), "bob_id", "Bob")
	_ = memStore.SaveMessage(context.Background(), "alice_id", store.Message{
		Sender:  "Bob",
		Time:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Payload: "Привет, это тестовое сообщение!",
	})

	appInstance := NewApp(memStore)
	ts := httptest.NewServer(http.HandlerFunc(appInstance.webhook))
	defer ts.Close()

	req := models.Request{
		Request: models.SimpleUtterance{
			Type:    models.TypeSimpleUtterance,
			Command: "Прочитай 1",
		},
		Session: models.Session{
			New:  false,
			User: models.RequestUser{UserID: "alice_id"},
		},
		Version: "1.0",
	}
	_ = sendPost(ts.URL, req)
	// В ответе будет текст сообщения
}

// Example_checkMessages демонстрирует проверку количества новых сообщений.
func Example_checkMessages() {
	memStore := &inMemoryStore{
		users:    make(map[string]string),
		messages: make(map[string][]store.Message),
	}
	_ = memStore.SaveMessage(context.Background(), "alice_id", store.Message{Sender: "Bob", Payload: "Привет!"})
	_ = memStore.SaveMessage(context.Background(), "alice_id", store.Message{Sender: "Charlie", Payload: "Как дела?"})

	appInstance := NewApp(memStore)
	ts := httptest.NewServer(http.HandlerFunc(appInstance.webhook))
	defer ts.Close()

	req := models.Request{
		Request: models.SimpleUtterance{
			Type:    models.TypeSimpleUtterance,
			Command: "что-то там",
		},
		Session: models.Session{
			New:  false,
			User: models.RequestUser{UserID: "alice_id"},
		},
		Version: "1.0",
	}
	_ = sendPost(ts.URL, req)
	// Ожидается ответ "Для вас 2 новых сообщений."
}

// Example_newSessionGreeting демонстрирует приветствие при новой сессии.
func Example_newSessionGreeting() {
	memStore := &inMemoryStore{
		users:    make(map[string]string),
		messages: make(map[string][]store.Message),
	}
	_ = memStore.SaveMessage(context.Background(), "alice_id", store.Message{Sender: "Bob", Payload: "Привет!"})

	appInstance := NewApp(memStore)
	ts := httptest.NewServer(http.HandlerFunc(appInstance.webhook))
	defer ts.Close()

	req := models.Request{
		Request: models.SimpleUtterance{
			Type:    models.TypeSimpleUtterance,
			Command: "что-то там",
		},
		Session: models.Session{
			New:  true,
			User: models.RequestUser{UserID: "alice_id"},
		},
		Timezone: "Europe/Moscow",
		Version:  "1.0",
	}
	_ = sendPost(ts.URL, req)
	// В ответе будет приветствие с текущим временем, затем количество сообщений
}