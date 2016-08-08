# PingRoll

Software to notify user on duty (through google calendar) the messages it receives from webhooks. It uses the Twilio API to send SMS.

## Configuration

* `PINGROLL_ICAL_URL` _(required)_: Google Calendar iCal URL
* `PINGROLL_RETHINK_URL` _(required)_: RethinkDB URL
* `PINGROLL_TWILIO_ACCOUNT_SID` _(required)_: Twilio Account SID
* `PINGROLL_TWILIO_AUTH_TOKEN` _(required)_: Twilio Auth Token
* `PINGROLL_TWILIO_FROM_NUMBER` _(required)_: Twilio From Number
* `PINGROLL_API_KEY` _(required)_: PingRoll API key

## Usage

Check `pingroll.apib` for the API blueprint. You register the users who may be on duty with the `/api/users` API. Then you register the webhook templates with the `/api/webhooks` API. You create events on your calendar with _summary_ reflecting the _name_ property of the users and pingroll updates its internal database with this calendar every hour. For each webhook you setup, the name is used to generate a unique friendly URL like `/webhook/:name`.

Every times a webhook receives a request, PingRoll transforms the request content into a `map[string]interface{}`. Then it applies the related template to this object. For more information about the templating engine see [text/template](https://godoc.org/text/template).

Then PingRoll looks in the database who's on duty and who should always receive notification. Finally for each recipient it makes a call to the Twilio API with the executed template as the body.

## TODO

- [ ] Protect webhook
- [ ] Multiple notification backend (such as mail and so on)
