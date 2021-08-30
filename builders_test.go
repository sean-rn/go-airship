package airship

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const templateIDA = "template-id-a"
const channelA = "channel-a"
const channelB = "channel-b"

func TestNewPushTemplatePayload(t *testing.T) {
	const expected = `{
		"audience": {
			"channel": ["channel-a"]
		},
		"device_types": [
			"ios", "android"
		],
		"merge_data": {
			"substitutions": {
				"AuthorFirstName": "The",
				"AuthorLastName": "Beat",
				"ActivityID": "2342"
			},
			"template_id": "template-id-a"
		}
	}`
	subs := map[string]string{
		"AuthorFirstName": "The",
		"AuthorLastName":  "Beat",
		"ActivityID":      "2342",
	}
	payload := MakePushTemplatePayload(templateIDA, []string{channelA}, subs)
	json, err := json.Marshal(&payload)
	require.Nil(t, err)
	assert.JSONEq(t, expected, string(json))
}

func TestNewSendPushPayload(t *testing.T) {
	const expected = `{
		"audience": {
			"channel": ["channel-a"]
		},
		"global_attributes": {
			"AuthorFirstName": "The",
			"AuthorLastName": "Beat",
			"ActivityID": "2402"
		},
		"notification": {
			"ios": {
				"template": {
					"template_id": "template-id-a"
				}
			},
			"android": {
				"template": {
					"template_id": "template-id-a"
				}
			}
		},
		"device_types": ["ios", "android"]
	}`
	subs := map[string]string{
		"AuthorFirstName": "The",
		"AuthorLastName":  "Beat",
		"ActivityID":      "2402",
	}
	payload := MakeSendPushPayload(templateIDA, []string{channelA}, subs)
	json, err := json.Marshal(&payload)
	require.Nil(t, err)
	assert.JSONEq(t, expected, string(json))
}

func TestSendPushPayload_WithActionAndExtra(t *testing.T) {
	const expected = `{
		"audience": {
			"channel": ["channel-a"]
		},
		"global_attributes": {
			"AuthorFirstName": "The",
			"AuthorLastName": "Beat",
			"ActivityID": "2402"
		},
		"notification": {
			"ios": {
				"template": {
					"template_id": "template-id-a"
				},
				"extra": {
					"shift_id": "12345"
				}
			},
			"android": {
				"template": {
					"template_id": "template-id-a"
				},
				"extra": {
					"shift_id": "12345"
				}
			},
			"actions": {
				"open": {
					"type": "url",
					"content": "https://xkcd.com/{{ActivityID}}"
				}
			}
		},
		"device_types": ["ios", "android"]
	}`
	subs := map[string]string{
		"AuthorFirstName": "The",
		"AuthorLastName":  "Beat",
		"ActivityID":      "2402",
	}
	payload := MakeSendPushPayload(templateIDA, []string{channelA}, subs,
		WithExtra(map[string]string{"shift_id": "12345"}),
		WithOpenURLAction("https://xkcd.com/{{ActivityID}}"))
	json, err := json.Marshal(&payload)
	require.Nil(t, err)
	assert.JSONEq(t, expected, string(json))
}

func TestNewSendPushPayload_WithDeepLink(t *testing.T) {
	const expected = `{
		"audience": {
			"channel": ["channel-a"]
		},
		"global_attributes": {
			"AuthorFirstName": "The",
			"AuthorLastName": "Beat",
			"ActivityID": "2402"
		},
		"notification": {
			"ios": {
				"template": {
					"template_id": "template-id-a"
				}
			},
			"android": {
				"template": {
					"template_id": "template-id-a"
				}
			},
			"actions": {
				"open": {
					"type": "deep_link",
					"content": "deep://link"
				}
			}
		},
		"device_types": ["ios", "android"]
	}`
	subs := map[string]string{
		"AuthorFirstName": "The",
		"AuthorLastName":  "Beat",
		"ActivityID":      "2402",
	}
	payload := MakeSendPushPayload(templateIDA, []string{channelA}, subs, WithDeepLinkAction("deep://link", ""))
	json, err := json.Marshal(&payload)
	require.Nil(t, err)
	assert.JSONEq(t, expected, string(json))
}
