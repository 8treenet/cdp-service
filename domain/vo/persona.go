package vo

import "encoding/json"

type ReqCreatePersona struct {
	Name      string          `json:"name"`
	Title     string          `json:"title"`
	DateRange int             `json:"dateRange"`
	XmlData   json.RawMessage `json:"xmlData"`
}
