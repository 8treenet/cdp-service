package vo

import "encoding/json"

type ReqCreateAnalysis struct {
	Name      string          `json:"name"`
	Title     string          `json:"title"`
	OutType   string          `json:"outType"`
	FeatureId int             `json:"featureId"`
	DateRange int             `json:"dateRange"`
	XmlData   json.RawMessage `json:"xmlData"`
}
