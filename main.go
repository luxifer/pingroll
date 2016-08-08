package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	r "gopkg.in/dancannon/gorethink.v2"
)

var (
	rc *r.Session
	tc *TwilioClient
)

func main() {
	router := echo.New()
	router.GET("/api/users", UserListHandler)
	router.PUT("/api/users", AddUserHandler)
	router.DELETE("/api/users/:id", RemoveUserHandler)
	router.GET("/api/webhooks", WebhookListHandler)
	router.PUT("/api/webhooks", AddWebhookHandler)
	router.DELETE("/api/webhooks/:id", RemoveWebhookHandler)
	router.POST("/webhook/:name", WebhookHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = ":1323"
	}

	icalURL := os.Getenv("PINGROLL_ICAL_URL")

	if icalURL == "" {
		log.Fatal("PINGROLL_ICAL_URL is required")
	}

	rehtinkUrl := os.Getenv("PINGROLL_RETHINK_URL")

	if rehtinkUrl == "" {
		log.Fatal("PINGROLL_RETHINK_URL is required")
	}

	twilioAccountSID := os.Getenv("PINGROLL_TWILIO_ACCOUNT_SID")

	if twilioAccountSID == "" {
		log.Fatal("PINGROLL_TWILIO_ACCOUNT_SID is required")
	}

	twilioAuthToken := os.Getenv("PINGROLL_TWILIO_AUTH_TOKEN")

	if twilioAuthToken == "" {
		log.Fatal("PINGROLL_TWILIO_AUTH_TOKEN is required")
	}

	twilioFromNumber := os.Getenv("PINGROLL_TWILIO_FROM_NUMBER")

	if twilioFromNumber == "" {
		log.Fatal("PINGROLL_TWILIO_FROM_NUMBER is required")
	}

	var err error
	rc, err = r.Connect(r.ConnectOpts{
		Address:  rehtinkUrl,
		Database: "pingroll",
	})

	if err != nil {
		log.Fatal(err)
	}

	tc = &TwilioClient{
		AccountSID: twilioAccountSID,
		AuthToken:  twilioAuthToken,
		FromNumber: twilioFromNumber,
	}

	ProcessIcal(icalURL)

	router.Run(standard.New(port))
}
