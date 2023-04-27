package fakemail

import (
	"sync"
	"time"
)

type MockEmailSender struct {
	// TODO add slog support

	emails        map[string]MockEmail
	lock          *sync.RWMutex
	lastUpdatedAt time.Time
}

func NewMockSender() *MockEmailSender {
	return &MockEmailSender{
		emails:        map[string]MockEmail{},
		lock:          &sync.RWMutex{},
		lastUpdatedAt: time.Now(),
	}
}

type MockEmail struct {
	From string   `json:"from"`
	To   []string `json:"to"`
	Body string   `json:"body"`

	Metadata map[string]string `json:"metadata"`
}

// Send(to []string, tplName TplID, tplVars map[string]interface{}) error
// BulkSend(templateID TplID, messages []PreparedEmail) error
