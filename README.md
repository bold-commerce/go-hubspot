[![Build Status](https://travis-ci.org/bold-commerce/go-hubspot.png)](https://travis-ci.org/bold-commerce/go-hubspot)

# go-hubspot
[Go](https://golang.org/) client for [HubSpot](https://app.hubspot.com)

*Note: This currently does not implement all HubSpot API endpoints*

## Install
```
go get github.com/bold-commerce/go-hubspot
```

## Unit Tests
To run the unit tests, install [ginkgo](https://onsi.github.io/ginkgo) and [gomega](https://onsi.github.io/gomega/) and run:

```
ginkgo -r
```

## Usage

```go
client := hubspot.NewClient("https://api.hubapi.com", "my-api-key")

// send single email
emailId := 12345678
err := client.SingleEmail(emailId, "tyler.durden@gmail.com")
if err != nil {
	log.Fatalf("hubspot error: %s", err.Error())
}
```
