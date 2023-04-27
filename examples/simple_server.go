package main

import (
	"log"
	"net/http"

	"github.com/cosmotek/fakemail"
)

func main() {
	mockSender := fakemail.NewMockSender()
	http.Handle("/sent_email/", mockSender.EmailViewer("/sent_email/"))

	log.Println("started email viewer at http://localhost:5050/sent_email")
	http.ListenAndServe(":5050", nil)
}
