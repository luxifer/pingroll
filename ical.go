package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/luxifer/ical"
)

func ProcessIcal(icalURL string) {
	FetchIcal(icalURL, time.Now())
	ticker := time.NewTicker(time.Hour) // update ical every hour
	go func() {
		for t := range ticker.C {
			FetchIcal(icalURL, t)
		}
	}()
}

func FetchIcal(icalURL string, tick time.Time) error {
	res, err := http.DefaultClient.Get(icalURL)

	if err != nil {
		log.Warn("Unable to fetch iCal")
		return err
	}

	calendar, err := ical.Parse(res.Body, nil)

	if err != nil {
		log.Warn("Unable to parse iCal")
		return err
	}

	for _, event := range calendar.Events {
		u, err := GetUserByName(event.Summary)

		if err != nil {
			return err
		}

		ts := &TimeSlot{
			UID:       event.UID,
			StartDate: event.StartDate,
			EndDate:   event.EndDate,
			User:      u,
		}
		err = UpdateTimeSlot(ts)

		if err != nil {
			return err
		}
	}

	log.Info("Ical updated")

	return nil
}
