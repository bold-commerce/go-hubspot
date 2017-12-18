package hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseUrl string
	apiKey  string
}

func NewClient(baseUrl, apiKey string) *Client {
	return &Client{
		baseUrl: baseUrl,
		apiKey:  apiKey,
	}
}

// Hubspot single send email API
// example: https://api.hubapi.com/email/public/v1/singleEmail/send?hapikey=demo
func (c *Client) SingleEmail(emailId int, emailTo string) error {

	req := SingleSendEmailRequest{
		EmailID: emailId,
		Message: Message{
			To: emailTo,
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("invalid request: ", err)
	}

	_, err = c.doRequest(request{
		URL:          fmt.Sprintf("%s/email/public/v1/singleEmail/send?hapikey=%s", c.baseUrl, c.apiKey),
		Method:       http.MethodPost,
		Body:         body,
		OkStatusCode: http.StatusOK,
	})

	return err
}

type request struct {
	URL          string
	Method       string
	Body         []byte
	OkStatusCode int
}

type response struct {
	Body       []byte
	StatusCode int
}

func (c *Client) doRequest(r request) (response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	if err != nil {
		return response{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response{}, err
	}

	if resp.StatusCode != r.OkStatusCode {
		return response{}, fmt.Errorf("Error: %s details: %s\n", resp.Status, body)
	}
	return response{Body: body, StatusCode: resp.StatusCode}, nil
}
