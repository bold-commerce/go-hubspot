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
	req := SendEmailRequest{
		EmailID: emailId,
		Message: Message{
			To: emailTo,
		},
	}
	return c.Email(req)
}

// Hubspot send email API
// example: https://api.hubapi.com/email/public/v1/singleEmail/send?hapikey=demo
func (c *Client) Email(emailRequest SendEmailRequest) error {

	body, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("invalid request: %s", err.Error())
	}
	_, err = c.doRequest(request{
		URL:          fmt.Sprintf("%s/email/public/v1/singleEmail/send?hapikey=%s", c.baseUrl, c.apiKey),
		Method:       http.MethodPost,
		Body:         body,
		OkStatusCode: http.StatusOK,
	})

	return err
}

// Hubspot Create or update a contact
// example https://api.hubapi.com/contacts/v1/contact/createOrUpdate/email/testingapis@hubspot.com/?hapikey=demo
func (c *Client) CreateOrUpdateContact(emailAddress string, properties []Property) (Response, error) {
	req := ContactBody{
		Properties: properties,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return Response{}, fmt.Errorf("invalid request: %s", err.Error())
	}

	response, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/contacts/v1/contact/createOrUpdate/email/%s/?hapikey=%s", c.baseUrl, emailAddress, c.apiKey),
		Method:       http.MethodPost,
		Body:         body,
		OkStatusCode: http.StatusOK,
	})

	return response, err
}

// Hubspot add contact to list
// example https://api.hubapi.com/contacts/v1/lists/226468/add?hapikey=demo
func (c *Client) AddContactsToList(emails []string, listId int) (Response, error) {
	return c.updateListWithContacts(listId, emails, "add")
}

func (c *Client) RemoveContactsFromList(emails []string, listId int) (Response, error) {
	return c.updateListWithContacts(listId, emails, "remove")
}

// Hubspot remove contact to list
// example https://api.hubapi.com/contacts/v1/lists/226468/remove?hapikey=demo
func (c *Client) updateListWithContacts(listId int, emails []string, method string) (Response, error) {
	req := ListBody{
		Emails: emails,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return Response{}, fmt.Errorf("invalid request: %s", err.Error())
	}

	response, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/contacts/v1/lists/%d/%s?hapikey=%s", c.baseUrl, listId, method, c.apiKey),
		Method:       http.MethodPost,
		Body:         body,
		OkStatusCode: http.StatusOK,
	})

	return response, err
}

// Hubspot add contact to workflow
// example https://api.hubapi.com/automation/v2/workflows/10900/enrollments/contacts/testingapis@hubspot.com?hapikey=demo
func (c *Client) AddContactToWorkFlow(email string, workflowId int) error {
	return c.updateWorkflowForClient(email, workflowId, http.MethodPost)
}

// Hubspot remove contact from workflow
// example https://api.hubapi.com/automation/v2/workflows/10900/enrollments/contacts/testingapis@hubspot.com?hapikey=demo
func (c *Client) RemoveContactFromWorkFlow(email string, workflowId int) error {
	return c.updateWorkflowForClient(email, workflowId, http.MethodDelete)
}

func (c *Client) updateWorkflowForClient(email string, workflowId int, method string) error {
	_, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/automation/v2/workflows/%d/enrollments/contacts/%s?hapikey=%s", c.baseUrl, workflowId, email, c.apiKey),
		Method:       method,
		OkStatusCode: http.StatusNoContent,
	})

	return err
}

type request struct {
	URL          string
	Method       string
	Body         []byte
	OkStatusCode int
}

type Response struct {
	Body       []byte
	StatusCode int
}

func (c *Client) doRequest(r request) (Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	if err != nil {
		return Response{}, err
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
