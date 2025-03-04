package api

import "net/url"

type URLRecordInfo struct {
	ID          Identificator `json:"id"`
	Name        string        `json:"name"`
	Destination string        `json:"destination"`
	Type        string        `json:"type"`
	ResourceURL url.URL       `json:"resource_url"`
	Delete      bool          `json:"delete"`
	Modify      bool          `json:"modify"`
}

type URLRecord struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Type        string `json:"type"`
}

func (c Client) CreateURLRecord(domain Domain, urlRecord URLRecord) (URLRecordInfo, error) {
	var result []URLRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(urlRecord).
		SetPathParam("domain", domain).
		Post("/dns/{domain}/url")

	if err := getResponseError(res); err != nil {
		return URLRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) GetURLRecords(domain Domain) ([]URLRecordInfo, error) {
	var result []URLRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		Get("/dns/{domain}/url")

	if err := getResponseError(res); err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetURLRecord(domain Domain, id Identificator) (URLRecordInfo, error) {
	var result []URLRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Get("/dns/{domain}/url/{id}")

	if err := getResponseError(res); err != nil {
		return URLRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) UpdateURLRecord(domain Domain, id Identificator, urlRecord URLRecord) (URLRecordInfo, error) {
	var result []URLRecordInfo

	res, _ := c.resty.R().
		SetResult(&result).
		SetContentType("application/json").
		SetBody(urlRecord).
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Put("/dns/{domain}/url/{id}")

	if err := getResponseError(res); err != nil {
		return URLRecordInfo{}, err
	}
	record := result[0]
	return record, nil
}

func (c Client) DeleteURLRecord(domain Domain, id Identificator) error {
	res, _ := c.resty.R().
		SetPathParam("domain", domain).
		SetPathParam("id", id).
		Delete("/dns/{domain}/url/{id}")

	return getResponseError(res)
}
