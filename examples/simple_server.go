package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cosmotek/fakemail"
)

func main() {
	mockSender := fakemail.NewMockSender()
	http.Handle("/sent_email/", mockSender.EmailViewer("/sent_email/"))

	for i := 0; i < 20; i++ {
		mockSender.Send(fakemail.MockEmail{
			From:    "seth@test.com",
			To:      []string{"bob@test.com"},
			Subject: "Hello World",
			Body:    fmt.Sprintf("<h1>Hello World %d</h1>", i),
			Metadata: map[string]string{
				"test": "123",
			},
		})
	}

	log.Println("started email viewer at http://localhost:5050/sent_email")
	http.ListenAndServe(":5050", nil)
}
