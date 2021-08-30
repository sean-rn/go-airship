package airship

// PushTemplatePayload https://docs.airship.com/api/ua/#schemas-pushtemplatepayload
type PushTemplatePayload struct {
	Audience    AudienceSelector `json:"audience" validate:"required"`
	DeviceTypes []string         `json:"device_types" validate:"required"` // "ios" and/or "android"
	MergeData   MergeData        `json:"merge_data" validate:"required"`
}

// AudienceSelector https://docs.airship.com/api/ua/#schemas-audienceselector
// Atomic Selector variant: https://docs.airship.com/api/ua/#schemas-atomicselector
type AudienceSelector struct {
	Channels   []string `json:"channel,omitempty"`
	NamedUsers []string `json:"named_user,omitempty"`
}

// MergeData is the merge_data field of a Push Template Payload
type MergeData struct {
	Substitutions map[string]string `json:"substitutions,omitempty"`
	TemplateID    string            `json:"template_id" validate:"required"`
}

// PushObject https://docs.airship.com/api/ua/#schemas-pushobject
type PushObject struct {
	Audience         AudienceSelector   `json:"audience" validate:"required"`
	DeviceTypes      interface{}        `json:"device_types" validate:"required"` // "all" or slice of "ios", "android", etc
	GlobalAttributes map[string]string  `json:"global_attributes,omitempty"`      // will be added to the global attributes rendering namespace for this push.
	Notification     NotificationObject `json:"notification"`                     // Probably yes required unless either message or in_app is present.
	// feed_references TODO - Probably don't need this
	// Message uaMessageCenterWithTemplate `json:"message"` // Probably not?  Either "Message Center Object" or "Message Center with Template"
}

// NotificationObject https://docs.airship.com/api/ua/#schemas-notificationobject
type NotificationObject struct {
	Alert   string                       `json:"alert,omitempty"`
	Actions *Actions                     `json:"actions,omitempty"`
	Android *AndroidOverrideWithTemplate `json:"android,omitempty"`
	IOS     *IOSOverrideWithTemplate     `json:"ios,omitempty"`
	Sms     *SMSOverrideWithTemplate     `json:"sms,omitempty"`
}

// SMSOverrideWithTemplate specifies an SMS message template to send.
// https://docs.airship.com/api/ua/#schemas-smsoverridewithtemplate
type SMSOverrideWithTemplate struct {
	Template     *TemplateRef `json:"template,omitempty"`
	ShortenLinks bool         `json:"shorten_links,omitempty"`
}

// AndroidOverrideWithTemplate https://docs.airship.com/api/ua/#schemas-androidoverridewithtemplate
type AndroidOverrideWithTemplate struct {
	Template    *TemplateRef      `json:"template,omitempty"`
	Actions     *Actions          `json:"actions,omitempty"`
	Extra       map[string]string `json:"extra,omitempty"`
	Sound       string            `json:"sound,omitempty"`
	CollapseKey string            `json:"collapse_key,omitempty"`
	Category    string            `json:"category,omitempty"`
	Title       string            `json:"title,omitempty"`
}

// IOSOverrideWithTemplate https://docs.airship.com/api/ua/#schemas-iosoverridewithtemplate
type IOSOverrideWithTemplate struct {
	Template   *TemplateRef      `json:"template,omitempty"`
	Actions    *Actions          `json:"actions,omitempty"`
	Extra      map[string]string `json:"extra,omitempty"`
	Sound      string            `json:"sound,omitempty"`
	Badge      int32             `json:"badge,omitempty"`
	CollapseID string            `json:"collapse_id,omitempty"`
	Category   string            `json:"category,omitempty"`
	Title      string            `json:"title,omitempty"`
}

// TemplateRef just holds a template ID under a key.
// One and only one of TemplateID and Fields may be populated
type TemplateRef struct {
	TemplateID string          `json:"template_id,omitempty" validate:"excluded_with=Fields"`
	Fields     *TemplateFields `json:"fields,omitempty" validate:"excluded_with=TemplateID"`
}

// TemplateFields allows specifying the template directly in the API call. Items in the field object are personalizable with handlebars.
type TemplateFields struct {
	Alert     string `json:"alert,omitempty"`
	Icon      string `json:"icon,omitempty"`
	IconColor string `json:"icon_color,omitempty"`
	Summary   string `json:"summary,omitempty"`
	Title     string `json:"title,omitempty"`
}

// Actions "Actions": Describes Actions to be performed by the SDK when a user interacts with the notification.
// https://docs.airship.com/api/ua/#schemas-actionsobject
type Actions struct {
	AddTag    []string      `json:"add_tag,omitempty"`
	RemoveTag []string      `json:"remove_tag,omitempty"`
	Share     string        `json:"string,omitempty"`
	Open      *uaOpenAction `json:"open,omitempty"`
}

// Action types fro the Action.Type field
const (
	ActionTypeDeepLink = "deep_link"
	ActionTypeOpenURL  = "url"
)

// uaOpenAction is the value of the "open" property of the "Actions" object.
type uaOpenAction struct {
	Type        string `json:"type" validate:"required"` // "url" or "deep_link"
	Content     string `json:"content"`                  // Used by URL and Deep Link
	FallbackURL string `json:"fallback_url,omitempty"`   // Used by Deep Link
}
