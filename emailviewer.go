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

const viewerHTML = `
<script type="module">
  import { createApp } from 'https://unpkg.com/petite-vue?module'
  
  createApp({
    emails: [],
	lastRefreshed: null,
	lastUpdated: null,
	paused: false,

	mounted() {
		this.getEmails();

		setInterval(() => {
			if (!this.paused) {
				this.getEmails();
			}
		}, 3000); 	
	},

	// methods
	getEmails() {
		fetch('%s/emails')
			.then((res) => res.json())
			.then((res) => {
				this.emails = res.emails;
				this.lastRefreshed = Date.now();
				this.lastUpdated = res.lastUpdatedAtUTC;
			}).catch((err) => console.warn('Something went wrong.', err));
	}
  }).mount()
</script>

<div v-scope @vue:mounted="mounted">
  <p>{{lastRefreshed}} {{lastUpdated}}</p>
  <button @click="paused = !paused">{{paused ? "unpause" : "pause"}} refresh</button>
  <div class="email" v-for="email in emails">
  	<h2>{{ email.subject }}</h2>
	<h2>{{ email.from }}</h2>
	<p>{{ email.body }}</p>
  </div>
</div>

<style>
.email {
	background-color: grey;
	margin: 5px;
	padding: 5px;
}
</style>
`

func (m *MockEmailSender) EmailViewer(pathPrefix string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		urlPath := strings.TrimPrefix(req.URL.Path, pathPrefix)

		if urlPath == "/" || urlPath == "" {
			io.WriteString(res, fmt.Sprintf(viewerHTML, pathPrefix))
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
