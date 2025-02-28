package api

import (
	"net/netip"
	"net/url"
)

type AAAARecordInfo struct {
	ID          Identificator `json:"id"`
	Name        string        `json:"name"`
	Destination netip.Addr    `json:"destination"`
	ResourceURL url.URL       `json:"resource_url"`
	Delete      bool          `json:"delete"`
	Modify      bool          `json:"modify"`
}

type AAAARecord struct {
	Name        string     `json:"name"`
	Destination netip.Addr `json:"destination"`
}

func (c Client) CreateAAAARecord(domain Domain, aaaaRecord AAAARecord) (AAAARecordInfo, error) {
	var result []AAAARecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(aaaaRecord).
		SetPathParam("domain", domain).
		Post("/dns/{domain}/aaaa")

	if err := getResponseError(res); err != nil {
		return AAAARecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) GetAAAARecords(domain Domain) ([]AAAARecordInfo, error) {
	var result []AAAARecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		Get("/dns/{domain}/aaaa")

	if err := getResponseError(res); err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetAAAARecord(domain Domain, id Identificator) (AAAARecordInfo, error) {
	var result []AAAARecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Get("/dns/{domain}/aaaa/{id}")

	if err := getResponseError(res); err != nil {
		return AAAARecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) UpdateAAAARecord(domain Domain, id Identificator, aaaaRecord AAAARecord) (AAAARecordInfo, error) {
	var result []AAAARecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(aaaaRecord).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Put("/dns/{domain}/aaaa/{id}")

	if err := getResponseError(res); err != nil {
		return AAAARecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) DeleteAAAARecord(domain Domain, id Identificator) error {
	res, _ := c.resty.R().
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Delete("/dns/{domain}/aaaa/{id}")

	return getResponseError(res)
}
