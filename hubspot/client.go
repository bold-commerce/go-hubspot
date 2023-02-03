package hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseUrl  string
	apiToken string
}

func NewClient(baseUrl, apiToken string) *Client {
	return &Client{
		baseUrl:  baseUrl,
		apiToken: apiToken,
	}
}

// Hubspot single send email API
// example: https://api.hubapi.com/email/public/v1/singleEmail/send
func (c *Client) SingleEmail(emailId int, emailTo string) error {
	req := SendEmailRequest{
		EmailID: emailId,
		Message: Message{
			To: emailTo,
		},
	}
	return c.Email(req)
}

// Hubspot send email API
// example: https://api.hubapi.com/email/public/v1/singleEmail/send
func (c *Client) Email(emailRequest SendEmailRequest) error {

	body, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("invalid request: %s", err.Error())
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiToken

	request := Request{
		URL:          fmt.Sprintf("%s/email/public/v1/singleEmail/send", c.baseUrl),
		Method:       http.MethodPost,
		Body:         body,
		OkStatusCode: http.StatusOK,
		Headers:      headers,
	}

	_, err = c.doRequest(request)

	return err
}

// Hubspot Create or update a contact
// example https://api.hubapi.com/contacts/v1/contact/createOrUpdate/email/testingapis@hubspot.com
func (c *Client) CreateOrUpdateContact(emailAddress string, properties []Property) (Response, error) {
	req := ContactBody{
		Properties: properties,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return Response{}, fmt.Errorf("invalid request: %s", err.Error())
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiToken

	request := Request{
		URL:          fmt.Sprintf("%s/contacts/v1/contact/createOrUpdate/email/%s", c.baseUrl, emailAddress),
		Method:       http.MethodPost,
		Body:         body,
		OkStatusCode: http.StatusOK,
		Headers:      headers,
	}

	response, err := c.doRequest(request)

	return response, err
}

// Hubspot add contact to list
// example https://api.hubapi.com/contacts/v1/lists/226468/add
func (c *Client) AddContactsToList(emails []string, listId int) (Response, error) {
	return c.updateListWithContacts(listId, emails, "add")
}

func (c *Client) RemoveContactsFromList(emails []string, listId int) (Response, error) {
	return c.updateListWithContacts(listId, emails, "remove")
}

// Hubspot remove contact to list
// example https://api.hubapi.com/contacts/v1/lists/226468/remove
func (c *Client) updateListWithContacts(listId int, emails []string, method string) (Response, error) {
	req := ListBody{
		Emails: emails,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return Response{}, fmt.Errorf("invalid request: %s", err.Error())
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiToken

	request := Request{
		URL:          fmt.Sprintf("%s/contacts/v1/lists/%d/%s", c.baseUrl, listId, method),
		Method:       http.MethodPost,
		Body:         body,
		OkStatusCode: http.StatusOK,
		Headers:      headers,
	}

	response, err := c.doRequest(request)

	return response, err
}

// Hubspot add contact to workflow
// example https://api.hubapi.com/automation/v2/workflows/10900/enrollments/contacts/testingapis@hubspot.com
func (c *Client) AddContactToWorkFlow(email string, workflowId int) error {
	return c.updateWorkflowForClient(email, workflowId, http.MethodPost)
}

// Hubspot remove contact from workflow
// example https://api.hubapi.com/automation/v2/workflows/10900/enrollments/contacts/testingapis@hubspot.com
func (c *Client) RemoveContactFromWorkFlow(email string, workflowId int) error {
	return c.updateWorkflowForClient(email, workflowId, http.MethodDelete)
}

func (c *Client) updateWorkflowForClient(email string, workflowId int, method string) error {

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiToken

	request := Request{
		URL:          fmt.Sprintf("%s/automation/v2/workflows/%d/enrollments/contacts/%s", c.baseUrl, workflowId, email),
		Method:       method,
		OkStatusCode: http.StatusNoContent,
		Headers:      headers,
	}

	_, err := c.doRequest(request)

	return err
}

type Request struct {
	URL          string
	Method       string
	Body         []byte
	OkStatusCode int
	Headers      map[string]string
}

type Response struct {
	Body       []byte
	StatusCode int
}

func (c *Client) doRequest(r Request) (Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	if err != nil {
		return Response{}, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	if resp.StatusCode != r.OkStatusCode {
		return Response{}, fmt.Errorf("Error: %s details: %s\n", resp.Status, body)
	}
	return Response{Body: body, StatusCode: resp.StatusCode}, nil
}
