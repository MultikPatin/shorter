package models

type AddedLink struct {
	CorrelationID string
	Short         string
	Origin        string
}

type OriginLink struct {
	CorrelationID string
	URL           string
}

type Result struct {
	CorrelationID string
	Result        string
}
