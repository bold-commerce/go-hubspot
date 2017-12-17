package hubspot

type Message struct {
	To  string   `json:"to"`
	Cc  []string `json:"cc"`
	Bcc []string `json:"bcc"`
}
type SingleSendEmailRequest struct {
	EmailID int     `json:"emailId"`
	Message Message `json:"message"`
}
