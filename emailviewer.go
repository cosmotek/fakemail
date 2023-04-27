package fakemail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type EmailsSendRequest struct {
	Emails []MockEmail `json:"emails"`
}

type EmailsGetRequest struct {
	Emails           map[string]MockEmail `json:"emails"`
	LastUpdatedAtUTC time.Time            `json:"lastUpdatedAtUTC"`
}

type EmailsSendResponse struct {
	EmailIDs []string `json:"emailIds"`
}

type EmailsDeleteRequest struct {
	DeleteAllEmails bool     `json:"deleteAllEmails"`
	DeleteEmailsIDs []string `json:"deleteEmailsIds"`
}

func (m *MockEmailSender) EmailViewer(pathPrefix string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		urlPath := strings.TrimPrefix(req.URL.Path, pathPrefix)

		if urlPath == "/" || urlPath == "" {
			io.WriteString(res, "Hello world")
			return
		}

		if urlPath != "emails" && urlPath != "emails/" {
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		}

		switch req.Method {
		case http.MethodGet:
			m.lock.RLock()
			defer m.lock.RUnlock()

			err := json.NewEncoder(res).Encode(EmailsGetRequest{
				Emails:           m.emails,
				LastUpdatedAtUTC: m.lastUpdatedAt.UTC(),
			})
			if err != nil {
				http.Error(res, fmt.Sprintf("Failed to encode emails response: %v", err), http.StatusInternalServerError)
				return
			}
			return

		case http.MethodPost:
			emailsSendReq := EmailsSendRequest{}
			err := json.NewDecoder(req.Body).Decode(&emailsSendReq)
			if err != nil {
				http.Error(res, fmt.Sprintf("Failed to decode send emails request body: %v", err), http.StatusBadRequest)
				return
			}

			m.lock.Lock()
			defer m.lock.Unlock()

			newIds := []string{}

			for _, mockEmail := range emailsSendReq.Emails {
				id := fmt.Sprintf("%d", time.Now().UTC().UnixNano())
				m.emails[id] = mockEmail
				newIds = append(newIds, id)
			}

			err = json.NewEncoder(res).Encode(EmailsSendResponse{
				EmailIDs: newIds,
			})
			if err != nil {
				http.Error(res, fmt.Sprintf("Failed to encode emails send response: %v", err), http.StatusInternalServerError)
				return
			}
			return

		case http.MethodDelete:
			deleteEmailsReq := EmailsDeleteRequest{}
			err := json.NewDecoder(req.Body).Decode(&deleteEmailsReq)
			if err != nil {
				http.Error(res, fmt.Sprintf("Failed to decode delete emails request body: %v", err), http.StatusBadRequest)
				return
			}

			m.lock.Lock()
			defer m.lock.Unlock()

			if deleteEmailsReq.DeleteAllEmails {
				m.emails = map[string]MockEmail{}
			} else {
				for _, id := range deleteEmailsReq.DeleteEmailsIDs {
					delete(m.emails, id)
				}
			}
			return

		default:
			http.Error(res, "Method not supported for this route", http.StatusNotFound)
			return
		}
	})
}
