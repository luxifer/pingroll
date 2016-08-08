# PingRoll

Software to notify user on duty (through google calendar) the messages it receives from webhooks. It uses the Twilio API to send SMS.

## Configuration

* `PINGROLL_ICAL_URL` _(required)_: Google Calendar iCal URL
* `PINGROLL_RETHINK_URL` _(required)_: RethinkDB URL
* `PINGROLL_TWILIO_ACCOUNT_SID` _(required)_: Twilio Account SID
* `PINGROLL_TWILIO_AUTH_TOKEN` _(required)_: Twilio Auth Token
* `PINGROLL_TWILIO_FROM_NUMBER` _(required)_: Twilio From Number
* `PINGROLL_API_KEY` _(required)_: PingRoll API key
