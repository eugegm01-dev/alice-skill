// Package models определяет структуры данных для взаимодействия с API Яндекс.Алисы.
package models

const (
    	// TypeSimpleUtterance — тип запроса для простой фразы.
	TypeSimpleUtterance = "SimpleUtterance"
)

// Request описывает запрос пользователя к навыку.
// Документация: https://yandex.ru/dev/dialogs/alice/doc/request.html
type Request struct {
	Request  SimpleUtterance `json:"request"`
	Session  Session         `json:"session"`
	Timezone string          `json:"timezone"`
	Version  string          `json:"version"`
}

// SimpleUtterance содержит команду пользователя.
type SimpleUtterance struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// Session описывает сессию пользователя
type Session struct {
	New  bool        `json:"new"`
	User RequestUser `json:"user"`
}

// RequestUser содержит идентификатор пользователя.
type RequestUser struct {
	UserID string `json:"user_id"`
}

// Response описывает ответ навыка.
// Документация: https://yandex.ru/dev/dialogs/alice/doc/response.html
type Response struct {
	Response ResponsePayload `json:"response"`
	Version  string          `json:"version"`
}

// ResponsePayload содержит текст ответа.
type ResponsePayload struct {
	Text string `json:"text"`
}
