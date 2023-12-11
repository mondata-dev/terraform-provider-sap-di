package sap_di

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetFactsheet - Returns specific factsheet
func (c *Client) GetFactsheet(connection string, dataset string) (*Factsheet, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/app/datahub-app-metadata/api/v1/catalog/connections/%s/datasets/%s/factsheet",
			c.HostURL,
			connection,
			dataset,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	factsheet := &Factsheet{}
	err = json.Unmarshal(body, factsheet)
	if err != nil {
		return nil, err
	}

	return factsheet, nil
}
