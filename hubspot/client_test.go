package hubspot_test

import (
	"net/http"
	"net/http/httptest"
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
})

const singleEmailResp = `{
   "sendResult":"SENT",
   "id":"62a4a958-0123-42f2-660e-5352d3b401ea",
   "eventId":{
      "id":"62a4a958-0123-42f2-660e-5352d3b401ea",
      "created":1513626117453
   }
}`
