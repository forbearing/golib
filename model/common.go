package model

type UserAgent struct {
	Source   string `json:"source" schema:"source"`
	Platform string `json:"platform" schema:"platform"`
	Engine   string `json:"engine" schema:"engine"`
	Browser  string `json:"browser" schema:"browser"`
}
