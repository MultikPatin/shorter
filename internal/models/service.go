package models

type AddedLink struct {
	CorrelationId string
	Short         string
	Origin        string
}

type OriginLink struct {
	CorrelationId string
	URL           string
}

type Result struct {
	CorrelationId string
	Result        string
}
