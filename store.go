package main

import (
	"strings"
	"time"

	r "gopkg.in/dancannon/gorethink.v2"
)

func ListUser() ([]*User, error) {
	res, err := r.Table("user").Run(rc)

	defer res.Close()

	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)

	if err := res.All(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func AddUser(u *User) error {
	res, err := r.Table("user").Insert(u).RunWrite(rc)

	if err != nil {
		return err
	}

	u.ID = res.GeneratedKeys[0]

	return nil
}

func RemoveUser(ID string) error {
	_, err := r.Table("user").Get(ID).Delete().RunWrite(rc)

	return err
}

func ListWebhook() ([]*Webhook, error) {
	res, err := r.Table("webhook").Run(rc)

	defer res.Close()

	if err != nil {
		return nil, err
	}

	webhooks := make([]*Webhook, 0)

	if err := res.All(&webhooks); err != nil {
		return nil, err
	}

	return webhooks, nil
}

func AddWebhook(w *Webhook) error {
	res, err := r.Table("webhook").Insert(w).RunWrite(rc)

	if err != nil {
		return err
	}

	w.ID = res.GeneratedKeys[0]

	return nil
}

func RemoveWebhook(ID string) error {
	_, err := r.Table("webhook").Get(ID).Delete().RunWrite(rc)

	return err
}

func GetWebhookByName(name string) (*Webhook, error) {
	res, err := r.Table("webhook").Filter(map[string]string{
		"Name": name,
	}).Limit(1).Run(rc)

	defer res.Close()

	if err != nil {
		return nil, err
	}

	var w Webhook

	if err := res.One(&w); err != nil {
		return nil, err
	}

	return &w, nil
}

func UpdateTimeSlot(ts *TimeSlot) error {
	res, err := r.Table("time_slot").Filter(map[string]string{
		"UID": ts.UID,
	}).Limit(1).Run(rc)

	defer res.Close()

	if err != nil {
		return err
	}

	if res.IsNil() {
		w, _ := r.Table("time_slot").Insert(ts).RunWrite(rc)
		ts.ID = w.GeneratedKeys[0]
		return nil
	}

	var tmpTs TimeSlot

	if err := res.One(&tmpTs); err != nil {
		return err
	}

	r.Table("time_slot").Get(tmpTs.ID).Update(ts).RunWrite(rc)
	ts.ID = tmpTs.ID

	return nil
}

func GetUserByName(name string) (*User, error) {
	res, err := r.Table("user").Filter(func(user r.Term) r.Term {
		return user.Field("Name").Downcase().Eq(strings.ToLower(name))
	}).Limit(1).Run(rc)

	defer res.Close()

	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	var u User

	if err := res.One(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func GetCurrentTimeSlot() ([]*TimeSlot, error) {
	now := time.Now()
	res, err := r.Table("time_slot").Filter(func(ts r.Term) r.Term {
		return ts.Field("StartDate").Le(now).And(ts.Field("EndDate").Gt(now))
	}).Run(rc)

	defer res.Close()

	if err != nil {
		return nil, err
	}

	var slots []*TimeSlot

	if err := res.All(&slots); err != nil {
		return nil, err
	}

	return slots, nil
}

func GetAlwaysUser() ([]*User, error) {
	res, err := r.Table("user").Filter(map[string]bool{
		"Always": true,
	}).Run(rc)

	defer res.Close()

	if err != nil {
		return nil, err
	}

	var users []*User

	if err := res.All(&users); err != nil {
		return nil, err
	}

	return users, nil
}
