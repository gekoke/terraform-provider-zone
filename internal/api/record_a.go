package api

import (
	"net/netip"
	"net/url"
)

type ARecordInfo struct {
	ID          Identificator `json:"id"`
	Name        string        `json:"name"`
	Destination netip.Addr    `json:"destination"`
	ResourceURL url.URL       `json:"resource_url"`
	Delete      bool          `json:"delete"`
	Modify      bool          `json:"modify"`
}

type ARecord struct {
	Name        string     `json:"name"`
	Destination netip.Addr `json:"destination"`
}

func (c Client) CreateARecord(domain Domain, aRecord ARecord) (ARecordInfo, error) {
	var result []ARecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(aRecord).
		SetPathParam("domain", domain).
		Post("/dns/{domain}/a")

	if err := getResponseError(res); err != nil {
		return ARecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) GetARecords(domain Domain) ([]ARecordInfo, error) {
	var result []ARecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		Get("/dns/{domain}/a")

	if err := getResponseError(res); err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetARecord(domain Domain, id Identificator) (ARecordInfo, error) {
	var result []ARecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Get("/dns/{domain}/a/{id}")

	if err := getResponseError(res); err != nil {
		return ARecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) UpdateARecord(domain Domain, id Identificator, aRecord ARecord) (ARecordInfo, error) {
	var result []ARecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(aRecord).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Put("/dns/{domain}/a/{id}")

	if err := getResponseError(res); err != nil {
		return ARecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) DeleteARecord(domain Domain, id Identificator) error {
	res, _ := c.resty.R().
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Delete("/dns/{domain}/a/{id}")

	return getResponseError(res)
}
