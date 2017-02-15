package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	ID     string `gorethink:"id,omitempty"`
	Name   string
	Number string
	Always bool
}

type Webhook struct {
	ID       string `gorethink:"id,omitempty"`
	Name     string
	Template string
}

type TimeSlot struct {
	ID        string `gorethink:"id,omitempty"`
	UID       string
	StartDate time.Time
	EndDate   time.Time
	User      *User
}

type Request struct {
	Method  string
	URL     string
	Header  http.Header
	Content []byte
}

type Event struct {
	ID      string `gorethink:"id,omitempty"`
	Request Request
	Webhook *Webhook
	User    *User
	Body    string
	SentAt  time.Time
}

type TwilioClient struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

func NewEvent(u *User, w *Webhook, body string, r *http.Request) *Event {
	content, _ := ioutil.ReadAll(r.Body)

	return &Event{
		Request: Request{
			Method:  r.Method,
			URL:     r.URL.String(),
			Header:  r.Header,
			Content: content,
		},
		Webhook: w,
		User:    u,
		Body:    body,
		SentAt:  time.Now(),
	}
}
