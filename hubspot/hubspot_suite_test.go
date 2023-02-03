package hubspot_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHubSpot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HubSpot Suite")
}
