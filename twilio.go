package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

func (tc *TwilioClient) SensSMS(to, body string) error {
	apiUrl := "https://api.twilio.com/2010-04-01/Accounts/" + tc.AccountSID + "/Messages.json"

	v := url.Values{}
	v.Set("From", tc.FromNumber)
	v.Set("To", to)
	v.Set("Body", body)

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(v.Encode()))

	if err != nil {
		return err
	}

	req.SetBasicAuth(tc.AccountSID, tc.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		log.WithFields(log.Fields{
			"status": resp.StatusCode,
			"to":     to,
		}).Error("Error while sending SMS")

		return errors.New("Unable to send SMS to " + to)
	}

	log.WithFields(log.Fields{
		"to": to,
	}).Info("SMS sent")

	return nil
}
