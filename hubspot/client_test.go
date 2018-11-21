package hubspot_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/bold-commerce/go-hubspot/hubspot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		client *hubspot.Client
	)

	Describe("New", func() {
		It("returns a new hubspot client", func() {
			client = hubspot.NewClient("https://api.hubapi.com", "my-api-key")
			Expect(client).ToNot(BeNil())
		})
	})

	Describe("SingleEmail", func() {
		var (
			server *httptest.Server
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				if strings.Contains(r.RequestURI, "/email/public/v1/singleEmail/send") {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(singleEmailResp))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())
		})

		AfterEach(func() {
			server.Close()
		})

		It("sends a single email", func() {
			err := client.SingleEmail(12345678, "tyler.durden@gmail.com")
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Describe("Email", func() {
		var (
			server *httptest.Server
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				if strings.Contains(r.RequestURI, "/email/public/v1/singleEmail/send") {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(singleEmailResp))
					b, err := ioutil.ReadAll(r.Body)
					Expect(err).NotTo(HaveOccurred())
					Expect(b).To(MatchJSON(b))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())
		})

		AfterEach(func() {
			server.Close()
		})

		It("sends a single customized email", func() {
			email := hubspot.SendEmailRequest{
				EmailID: 2853049635,
				Message: hubspot.Message{
					To:     "example@hubspot.com",
					SendID: "foobar",
				},
				ContactProperties: []hubspot.MergeField{
					{Name: "first_name", Value: "John"},
				},
				CustomProperties: []hubspot.MergeField{
					{Name: "item_1", Value: "something they bought"},
				}}
			err := client.Email(email)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Describe("CreateOrUpdateContact", func() {
		var (
			server *httptest.Server
		)

		BeforeEach(func() {
			response, _ := json.Marshal(map[string]interface{}{"vid": 751, "isNew": true})

			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				if strings.Contains(r.RequestURI, "/contacts/v1/contact/createOrUpdate/email/") {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())
		})

		AfterEach(func() {
			server.Close()
		})

		It("creates an email an returns response from api", func() {
			properties := []hubspot.Property{hubspot.Property{}}
			response, err := client.CreateOrUpdateContact("gord.currie@boldcommerce.com", properties)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(response.Body)).To(Equal(`{"isNew":true,"vid":751}`))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Describe("AddContactsToList", func() {
		var (
			server *httptest.Server
		)

		BeforeEach(func() {
			response, _ := json.Marshal(map[string]interface{}{"updated": []int{751}, "discarded": []int{}, "invalidVids": []int{}, "invalidEmails": []string{}})

			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				match, _ := regexp.MatchString(".*contacts/v1/lists/\\d*/add\\?hapikey=.*", r.RequestURI)
				if match {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())
		})

		AfterEach(func() {
			server.Close()
		})

		It("removes emails from list and returns response from api", func() {
			emails := []string{"gord.currie@boldcommerce.com"}
			response, err := client.AddContactsToList(emails, 5)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Describe("RemoveContactsFromList", func() {
		var (
			server *httptest.Server
		)

		BeforeEach(func() {
			response, _ := json.Marshal(map[string]interface{}{"updated": []int{751}, "discarded": []int{}, "invalidVids": []int{}, "invalidEmails": []string{}})

			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				match, _ := regexp.MatchString(".*contacts/v1/lists/\\d*/remove\\?hapikey=.*", r.RequestURI)
				if match {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())
		})

		AfterEach(func() {
			server.Close()
		})

		It("removes emails from list and returns response from api", func() {
			emails := []string{"gord.currie@boldcommerce.com"}
			response, err := client.RemoveContactsFromList(emails, 5)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Describe("AddContactToWorkFlow", func() {
		var (
			server *httptest.Server
		)

		AfterEach(func() {
			server.Close()
		})

		It("removes email from workflow", func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				match, _ := regexp.MatchString(".*automation/v2/workflows/\\d*/enrollments/contacts/.*", r.RequestURI)
				if match {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())

			err := client.AddContactToWorkFlow("gord.currie@boldcommerce.com", 2494115)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error if the contact email is not found", func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				match, _ := regexp.MatchString(".*automation/v2/workflows/\\d*/enrollments/contacts/.*", r.RequestURI)
				if match {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())

			err := client.AddContactToWorkFlow("gord.currie@boldcommerce.com", 2494115)
			Expect(err).To(MatchError(ContainSubstring("Error: 404 Not Found")))
		})
	})

	Describe("RemoveContactFromWorkFlow", func() {
		var (
			server *httptest.Server
		)

		AfterEach(func() {
			server.Close()
		})

		It("removes email from workflow", func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				match, _ := regexp.MatchString(".*automation/v2/workflows/\\d*/enrollments/contacts/.*", r.RequestURI)
				if match {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())

			err := client.RemoveContactFromWorkFlow("gord.currie@boldcommerce.com", 2494115)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error if the contact email is not found", func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				match, _ := regexp.MatchString(".*automation/v2/workflows/\\d*/enrollments/contacts/.*", r.RequestURI)
				if match {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			client = hubspot.NewClient(server.URL, "my-api-key")
			Expect(client).ToNot(BeNil())

			err := client.RemoveContactFromWorkFlow("gord.currie@boldcommerce.com", 2494115)
			Expect(err).To(MatchError(ContainSubstring("Error: 404 Not Found")))
		})
	})
})

const singleEmailResp = `{
   "sendResult":"SENT",
   "id":"62a4a958-0123-42f2-660e-5352d3b401ea",
   "eventId":{
      "id":"62a4a958-0123-42f2-660e-5352d3b401ea",
      "created":1513626117453
   }
}`
const emailRequest = `
{
    "emailId": 2853049635,
    "message": {
        "to": "example@hubspot.com",
        "sendId": "foobar"
    },
    "contactProperties": [
        {
            "name": "first_name",
            "value": "John"
        }
    ],
    "customProperties": [
        {
            "name": "item_1",
            "value": "something they bought"
        }
    ]
}
`
