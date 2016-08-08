package main

import (
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

type TwilioClient struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}
