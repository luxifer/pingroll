package main

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

func UserListHandler(c echo.Context) error {
	users, err := ListUser()

	if err != nil {
		log.Warn(err)
		return err
	}

	collection := make(map[string]interface{})
	collection["data"] = users

	return c.JSON(http.StatusOK, collection)
}

func AddUserHandler(c echo.Context) error {
	u := new(User)
	u.Always = false // By default, do not send SMS when no TimeSlot

	if err := c.Bind(u); err != nil {
		log.Warn(err)
		return err
	}

	if err := AddUser(u); err != nil {
		log.Warn(err)
		return err
	}

	log.WithFields(log.Fields{
		"user": u.ID,
	}).Info("New user added")

	return c.JSON(http.StatusCreated, u)
}

func RemoveUserHandler(c echo.Context) error {
	ID := c.Param("id")

	if err := RemoveUser(ID); err != nil {
		log.Warn(err)
		return err
	}

	log.WithFields(log.Fields{
		"user": ID,
	}).Info("User removed")

	return c.NoContent(http.StatusNoContent)
}

func WebhookListHandler(c echo.Context) error {
	webhooks, err := ListWebhook()

	if err != nil {
		log.Warn(err)
		return err
	}

	collection := make(map[string]interface{})
	collection["data"] = webhooks

	return c.JSON(http.StatusOK, collection)
}

func AddWebhookHandler(c echo.Context) error {
	w := new(Webhook)

	if err := c.Bind(w); err != nil {
		log.Warn(err)
		return err
	}

	if err := AddWebhook(w); err != nil {
		log.Warn(err)
		return err
	}

	log.WithFields(log.Fields{
		"webhook": w.ID,
	}).Info("New webhook added")

	return c.JSON(http.StatusCreated, w)
}

func RemoveWebhookHandler(c echo.Context) error {
	ID := c.Param("id")

	if err := RemoveWebhook(ID); err != nil {
		log.Warn(err)
		return err
	}

	log.WithFields(log.Fields{
		"webhook": ID,
	}).Info("Webhook removed")

	return c.NoContent(http.StatusNoContent)
}

func WebhookHandler(c echo.Context) error {
	name := c.Param("name")
	w, err := GetWebhookByName(name)

	if err != nil {
		log.Warn(err)
		return echo.NewHTTPError(http.StatusNotFound)
	}

	t := template.Must(template.New(name).Parse(w.Template))
	var data map[string]interface{}

	ct := c.Request().Header().Get(echo.HeaderContentType)

	switch ct {
	case echo.MIMEApplicationForm, echo.MIMEMultipartForm:
		data = make(map[string]interface{})
		for key, value := range c.FormParams() {
			if len(value) == 1 {
				data[key] = value[0]
			} else {
				data[key] = value
			}
		}
	default:
		if err := c.Bind(&data); err != nil {
			log.Warn(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}

	var b bytes.Buffer
	err = t.Execute(&b, data)

	if err != nil {
		log.Warn(err)
		return err
	}

	message := b.String()
	message = fmt.Sprintf("%s:\n%s", name, message)

	slots, err := GetCurrentTimeSlot()
	if err != nil {
		log.Warn(err)
		return err
	}

	users, err := GetAlwaysUser()
	if err != nil {
		log.Warn(err)
		return err
	}

	receivers := make(map[string]*User)

	// execute for user on duty
	for _, ts := range slots {
		receivers[ts.User.ID] = ts.User
	}

	// only send to always if nobody's on duty
	if len(receivers) == 0 {
		// execute for user with always on
		for _, u := range users {
			receivers[u.ID] = u
		}
	}

	// loop over receivers to prevent recipient message duplication
	for _, u := range receivers {
		log.WithFields(log.Fields{
			"user":    u.Name,
			"webhook": name,
		}).Info("Webhook triggered")
		tc.SensSMS(u.Number, message)
	}

	return c.NoContent(http.StatusCreated)
}
