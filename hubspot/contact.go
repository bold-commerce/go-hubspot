package hubspot

type ContactBody struct {
	Properties []Property `json:"properties"`
}

type Property struct {
	Property string      `json:"property"`
	Value    interface{} `json:"value"`
}

type ListBody struct {
	Emails []string `json:"emails"`
}
