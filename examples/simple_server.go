package main

import (
	"log"
	"net/http"

	"github.com/cosmotek/fakemail"
)

func main() {
	mockSender := fakemail.NewMockSender()
	http.Handle("/sent_email/", mockSender.EmailViewer("/sent_email/"))

	mockSender.Send(fakemail.MockEmail{
		From:    "seth@test.com",
		To:      []string{"bob@test.com"},
		Subject: "Hello World",
		Body:    "<h1>Hello World</h1>",
		Metadata: map[string]string{
			"test": "123",
		},
	})

	mockSender.Send(fakemail.MockEmail{
		From:    "seth@test.com",
		To:      []string{"bob@test.com"},
		Subject: "Hello World",
		Body:    "<h1>Hello World</h1>",
		Metadata: map[string]string{
			"test": "123",
		},
	})

	mockSender.Send(fakemail.MockEmail{
		From:    "seth@test.com",
		To:      []string{"bob@test.com"},
		Subject: "Hello World",
		Body:    "<h1>Hello World</h1>",
		Metadata: map[string]string{
			"test": "123",
		},
	})

	log.Println("started email viewer at http://localhost:5050/sent_email")
	http.ListenAndServe(":5050", nil)
}
