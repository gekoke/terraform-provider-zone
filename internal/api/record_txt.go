package api

import "net/url"

type TXTRecordInfo struct {
	ID          Identificator `json:"id"`
	Name        string        `json:"name"`
	Destination string        `json:"destination"`
	ResourceURL url.URL       `json:"resource_url"`
	Delete      bool          `json:"delete"`
	Modify      bool          `json:"modify"`
}

type TXTRecord struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
}

func (c Client) CreateTXTRecord(domain Domain, txtRecord TXTRecord) (TXTRecordInfo, error) {
	var result []TXTRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(txtRecord).
		SetPathParam("domain", domain).
		Post("/dns/{domain}/txt")

	if err := getResponseError(res); err != nil {
		return TXTRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) GetTXTRecords(domain Domain) ([]TXTRecordInfo, error) {
	var result []TXTRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		Get("/dns/{domain}/txt")

	if err := getResponseError(res); err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetTXTRecord(domain Domain, id Identificator) (TXTRecordInfo, error) {
	var result []TXTRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Get("/dns/{domain}/txt/{id}")

	if err := getResponseError(res); err != nil {
		return TXTRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) UpdateTXTRecord(domain Domain, id Identificator, txtRecord TXTRecord) (TXTRecordInfo, error) {
	var result []TXTRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(txtRecord).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Put("/dns/{domain}/txt/{id}")

	if err := getResponseError(res); err != nil {
		return TXTRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) DeleteTXTRecord(domain Domain, id Identificator) error {
	res, _ := c.resty.R().
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Delete("/dns/{domain}/txt/{id}")

	return getResponseError(res)
}
