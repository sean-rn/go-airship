package airship

import (
	"encoding/json"
	"fmt"
	"time"
)

// CreateAndSend is a create-and-send request body for Urban Airship
// https://docs.airship.com/api/ua/#operation-api-create-and-send-post
// https://docs.airship.com/api/ua/#schemas-sms for the SMS varient
type CreateAndSend struct {
	Audience     createAndSendAudience `json:"audience"`
	Notification NotificationObject    `json:"notification"`
	DeviceTypes  []string              `json:"device_types"`
}

// CreateAndSendSMSTarget defines an audience target for where to send an SMS message
type CreateAndSendSMSTarget struct {
	MSISDN  string    `json:"ua_msisdn"`   // The phone number of a mobile device.
	OptedIn time.Time `json:"ua_opted_in"` // The date/time when the user (msisdn) opted in to messages from the sender
	Sender  string    `json:"ua_sender"`   // The long or short code your SMS messages are sent from.
}

// Intermediate struct to create the necessary wrapper object with "create_and_send" property around the acutal array.
type createAndSendAudience struct {
	CreateAndSend []createAndSendAudienceEntry `json:"create_and_send"`
}

// Entries in CreateAndSend audience combine the target channel info and template substitutions into one object.
type createAndSendAudienceEntry struct {
	target        interface{} // CreateAndSendSMSTarget
	substitutions map[string]string
}

// MarshalJSON is overridden in order to merge the extra substitution fields in with the standard fields as is required by API.
func (s createAndSendAudienceEntry) MarshalJSON() ([]byte, error) {
	sBytes, err := json.Marshal(s.target) // Serialize the built-in fields
	if err != nil {
		return nil, err
	}
	if s.substitutions == nil || len(s.substitutions) == 0 {
		return sBytes, nil // The map is empty, so just use the struct's marshalling
	}
	mBytes, err := json.Marshal(s.substitutions)
	if err != nil {
		return nil, err
	}
	if len(sBytes) <= 2 {
		return mBytes, nil // The struct was empty (shouldn't happen but ok...)
	}
	sBytes[len(sBytes)-1] = ','               // Overwrite final '}' with ','
	return append(sBytes, mBytes[1:]...), nil // Append the serialization of the map
}

// MakeCreateAndSendSMSPayload creates a new create-and-send template payload to
// send an SMS message to the recipients using a message template already configured in Airship.
func MakeCreateAndSendSMSPayload(templateID string, subs map[string]string, shortenLinks bool, targets []CreateAndSendSMSTarget) (*CreateAndSend, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("airship: must specify at least one SMS destination")
	}
	audEntries := make([]createAndSendAudienceEntry, len(targets))
	for i := range targets {
		audEntries[i] = createAndSendAudienceEntry{
			target:        targets[i],
			substitutions: subs,
		}
	}

	return &CreateAndSend{
		Audience:    createAndSendAudience{CreateAndSend: audEntries},
		DeviceTypes: []string{"sms"},
		Notification: NotificationObject{
			Sms: &SMSOverrideWithTemplate{
				Template:     &TemplateRef{TemplateID: templateID},
				ShortenLinks: shortenLinks,
			},
		},
	}, nil
}
