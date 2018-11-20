package hubspot

type Message struct {
	To          string   `json:"to"`
	From        string   `json:"from,omitempty"`
	SendID      string   `json:"sendId,omitempty"`
	ReplyTo     string   `json:"replyTo,omitempty"`
	ReplyToList []string `json:"replyToList,omitempty"`
	Cc          []string `json:"cc,omitempty"`
	Bcc         []string `json:"bcc,omitempty"`
}
type MergeField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type SendEmailRequest struct {
	EmailID           int          `json:"emailId"`
	Message           Message      `json:"message"`
	ContactProperties []MergeField `json:"contactProperties,omitempty"`
	CustomProperties  []MergeField `json:"customProperties,omitempty"`
}
