package airship

// PushToTemplate invokes the Airship "Push to Template" API
// https://docs.airship.com/api/ua/#operation-api-templates-push-post
func PushToTemplate(templateID string, channels []string, substitutions map[string]string) {
	body := MakePushTemplatePayload(templateID, channels, substitutions)
	_ = body
	// return cfg.invokeEndpoint(http.Post(), "/api/templates/push", &ua)
}

// SendPush invokes the Airship "Send a Push" API
// https://docs.airship.com/api/ua/#operation-api-push-post
func SendPush(templateID string, channels []string, substitutions map[string]string) {
	body := MakeSendPushPayload(templateID, channels, substitutions)
	_ = body
	// return cfg.invokeEndpoint(http.Post(), "/api/push", &ua)
}
