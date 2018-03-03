package modules

import (
	"gopkg.in/resty.v1"
	"fmt"
	//"encoding/json"
	"os"
)

type FourSquareApi struct {
	Request *resty.Request
}

type FourSquareResponse struct {

}

func (api *FourSquareApi) InitClient() {
	api.Request = resty.R().SetQueryParams(map[string]string{
		"client_id": os.Getenv("API_FOURSQUARE_CLIENT_ID"),
		"client_secret": os.Getenv("API_FOURSQUARE_CLIENT_SECRET"),
		"v" : "20170801",
	})
}


func (api *FourSquareApi) GetVenues(query string) string {
	resp, err := api.Request.
		SetQueryParams(map[string]string{
		"ll": "59.93,30.31",
		"query" : query,
	}).
		SetHeader("Accept", "application/json").
		Get("https://api.foursquare.com/v2/venues/search")

	if err != nil {
		fmt.Println(err)
	}

	//var s ApiResponse
	//
	//json.Unmarshal(resp.Body(), &s)

	return  resp.String()
}
