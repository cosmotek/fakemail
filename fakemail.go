package fakemail

import (
	"fmt"
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

	Subject string `json:"subject"`
	Body    string `json:"body"`

	Metadata map[string]string `json:"metadata"`
}

func (m *MockEmailSender) Send(emails ...MockEmail) []string {
	m.lock.Lock()
	defer m.lock.Unlock()

	ids := []string{}
	for _, email := range emails {
		id := fmt.Sprintf("%d", time.Now().UTC().UnixNano())
		m.emails[id] = email
	}

	return ids
}
