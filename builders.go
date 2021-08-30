package airship

// MakePushTemplatePayload creates a new push template payload
func MakePushTemplatePayload(templateID string, channels []string, substitutions map[string]string) PushTemplatePayload {
	return PushTemplatePayload{
		Audience:    AudienceSelector{Channels: channels},
		DeviceTypes: []string{"ios", "android"},
		MergeData: MergeData{
			Substitutions: substitutions,
			TemplateID:    templateID,
		},
	}
}

// MakeSendPushPayload creates a new push template payload
func MakeSendPushPayload(templateID string, channels []string, substitutions map[string]string, options ...PushNotificationOption) PushObject {
	notification := NotificationObject{
		Android: &AndroidOverrideWithTemplate{
			Template: &TemplateRef{TemplateID: templateID},
		},
		IOS: &IOSOverrideWithTemplate{
			Template: &TemplateRef{TemplateID: templateID},
		},
	}
	for _, fn := range options {
		fn(&notification)
	}
	return PushObject{
		Audience:         AudienceSelector{Channels: channels},
		DeviceTypes:      []string{"ios", "android"},
		GlobalAttributes: substitutions,
		Notification:     notification,
	}
}

//
// Mutator Config -- Experimental
//

// PushNotificationOption is a mutator-function based option for sending a push notification.
type PushNotificationOption = func(notif *NotificationObject)

// WithDeepLinkAction adds an Open Deep Link action to the notification.
func WithDeepLinkAction(deepURL, fallbackURL string) PushNotificationOption {
	return func(notif *NotificationObject) {
		if notif.Actions == nil {
			notif.Actions = &Actions{}
		}
		notif.Actions.Open = &uaOpenAction{
			Type:        ActionTypeDeepLink,
			Content:     deepURL,
			FallbackURL: fallbackURL,
		}
	}
}

// WithOpenURLAction adds an Open URL action to the notification.
func WithOpenURLAction(url string) PushNotificationOption {
	return func(notif *NotificationObject) {
		if notif.Actions == nil {
			notif.Actions = &Actions{}
		}
		notif.Actions.Open = &uaOpenAction{
			Type:    ActionTypeOpenURL,
			Content: url,
		}
	}
}

// WithExtra adds "extra" data to IOS & Android notification overrides.
func WithExtra(extra map[string]string) PushNotificationOption {
	return func(notif *NotificationObject) {
		if notif.Android != nil {
			notif.Android.Extra = extra
		}
		if notif.IOS != nil {
			notif.IOS.Extra = extra
		}
	}
}

// WithShortenLinks sets the ShortenLinks value on the SMS notification overrides.
func WithShortenLinks(value bool) PushNotificationOption {
	return func(notif *NotificationObject) {
		if notif.Sms != nil {
			notif.Sms.ShortenLinks = value
		}
	}
}
