package modules

import (
	"gopkg.in/resty.v1"
	"fmt"
	"encoding/json"
)

type GisApi struct {

}

type ApiResponse struct {
	Result Result `json:"result"`
}

type Result struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name string `json:"name"`
	Ads Ads `json:"ads"`
}

type Ads struct {
	Link Link `json:"link"`
}

type Link struct {
	Value string `json:"value"`
	Text string `json:"text"`
}

func (api GisApi) GetItems() ApiResponse {
	resp, err := resty.R().
		SetQueryParams(map[string]string{
		"key": "ruhmxf0953",
		"q" : "поесть",
		"page": "1",
		"page_size": "50",
		"region_id": "38",
	}).
		SetHeader("Accept", "application/json").
		Get("https://catalog.api.2gis.ru/3.0/items")

	if err != nil {
		fmt.Println(err)
	}

	var s ApiResponse

	json.Unmarshal(resp.Body(), &s)

	return  s
}
