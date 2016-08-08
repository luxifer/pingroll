package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	r "gopkg.in/dancannon/gorethink.v2"
)

var (
	rc     *r.Session
	tc     *TwilioClient
	apiKey string
)

func main() {
	router := echo.New()
	api := router.Group("/api")
	api.Use(middleware.BasicAuth(AuthMiddleware))
	api.GET("/users", UserListHandler)
	api.PUT("/users", AddUserHandler)
	api.DELETE("/users/:id", RemoveUserHandler)
	api.GET("/webhooks", WebhookListHandler)
	api.PUT("/webhooks", AddWebhookHandler)
	api.DELETE("/webhooks/:id", RemoveWebhookHandler)
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

	apiKey = os.Getenv("PINGROLL_API_KEY")

	if apiKey == "" {
		log.Fatal("PINGROLL_API_KEY is required")
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
