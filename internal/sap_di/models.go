package sap_di

type Factsheet struct {
	Metadata FactsheetMetadata `json:"metadata"`
	Columns  []FactsheetColumn `json:"columns"`
}

type FactsheetMetadata struct {
	Name         string                 `json:"name"`
	Uri          string                 `json:"uri"`
	ConnectionId string                 `json:"connectionId"`
	Descriptions []FactsheetDescription `json:"descriptions"`
}

type FactsheetColumn struct {
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Descriptions []FactsheetDescription `json:"descriptions"`
}

type FactsheetDescription struct {
	Origin string `json:"origin"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}
