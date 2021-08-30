package airship

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newCreateAndSendPayload creates a new create-and-send template payload
func TestNewCreateAndSendSMSPayload(t *testing.T) {
	const expected = `{
		"audience": {
			"create_and_send": [
				{
					"ua_msisdn": "19785551212",
					"ua_opted_in": "2021-03-27T20:07:43Z",
					"ua_sender": "12062071886",
					"AuthorFirstName": "The",
					"AuthorLastName": "Beat",
					"ActivityID": "2342"
				}
			]
		},
		"device_types": [
			"sms"
		],
		"notification": {
			"sms": {
				"template": {
					"template_id": "template-id-a"
				},
				"shorten_links": true
			}
		}
	}`

	subs := map[string]string{
		"AuthorFirstName": "The",
		"AuthorLastName":  "Beat",
		"ActivityID":      "2342",
	}
	target := CreateAndSendSMSTarget{
		MSISDN:  "19785551212",
		OptedIn: time.Date(2021, 3, 27, 20, 7, 43, 0, time.UTC).UTC(),
		Sender:  "12062071886",
	}

	payload, err := MakeCreateAndSendSMSPayload("template-id-a", subs, true, []CreateAndSendSMSTarget{target})
	require.Nil(t, err)
	json, err := json.Marshal(&payload)
	require.Nil(t, err)
	assert.JSONEq(t, expected, string(json))
}

func TestCreateAndSendAudienceEntry_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		input    createAndSendAudienceEntry
		expected string
	}{
		{
			name: "with nil substitutions",
			input: createAndSendAudienceEntry{
				target: CreateAndSendSMSTarget{
					MSISDN:  "19785551212",
					OptedIn: time.Date(2021, 3, 27, 20, 7, 43, 0, time.UTC).UTC(),
					Sender:  "12062071886",
				},
				substitutions: nil,
			},
			expected: `{
				"ua_msisdn": "19785551212",
				"ua_opted_in": "2021-03-27T20:07:43Z",
				"ua_sender": "12062071886"
			  }`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.input)
			require.Nil(t, err)
			assert.JSONEq(t, tt.expected, string(bytes))
		})
	}
}

// func TestSendTemplatedSMSMessage(t *testing.T) {
// 	assert := assert.New(t)

// 	expectedBody := `{
// 		"audience": {
// 		  "create_and_send": [
// 			{
// 			  "ua_msisdn": "19785551212",
// 			  "ua_opted_in": "2021-03-27T20:07:43Z",
// 			  "ua_sender": "12062071886",
// 			  "ShiftDate": "2/24",
// 			  "ShiftID": "1942",
// 			  "first_name": "Testy McTesterson"
// 			}
// 		  ]
// 		},
// 		"device_types": [ "sms" ],
// 		"notification": {
// 		  "sms": {
// 			"shorten_links": true,
// 			"template": {
// 			  "template_id": "test-template-id"
// 			}
// 		  }
// 		}
// 	  }`

// 	client := httpmock.NewHandlerClient(func(rw http.ResponseWriter, req *http.Request) {
// 		// Validate request parameters and body
// 		assert.Equal("POST", req.Method)
// 		assert.Equal("https://go.urbanairship.com/api/create-and-send", req.URL.String())
// 		assert.Equal("application/json", req.Header.Get("Content-Type"))
// 		assert.Equal("Bearer test-ua-token", req.Header.Get("Authorization"))
// 		assert.Equal("application/vnd.urbanairship+json; version=3;", req.Header.Get("Accept"))
// 		assertBodyJSONEqual(t, expectedBody, req.Body)
// 		// Write response
// 		rw.Write([]byte(`{"ok": true,"operation_id": "df6a6b50","push_ids": ["9d78a53b"],"message_ids": [], "content_urls": []}`))
// 	})

// 	testConnection := New(WithHTTPClient(client), WithBearerAuth(TestBearerToken))

// 	// Test Data
// 	subs := map[string]string{
// 		"ShiftDate":  "2/24",
// 		"ShiftID":    "1942",
// 		"first_name": "Testy McTesterson",
// 	}
// 	target := CreateAndSendSMSTarget{
// 		MSISDN:  "19785551212",
// 		OptedIn: time.Date(2021, 3, 27, 20, 7, 43, 0, time.UTC).UTC(),
// 		Sender:  "12062071886",
// 	}

// 	// Invoke!
// 	err := testConnection.SendTemplatedSMSMessage("test-template-id", subs, target)
// 	require.Nil(t, err)
// }
