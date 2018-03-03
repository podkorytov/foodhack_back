package modules

import (
	"gopkg.in/resty.v1"
	"fmt"
	"encoding/json"
	"os"
)

type FourSquareApi struct {
	Request *resty.Request
}

type FourSquareResponse struct {
	Response Response `json:"response"`
}

type PhotoResponse struct {
	Photos Photos `json:"photos"`
}

type Photos struct {
	Items []Photo `json:"items"`
}

type FourSquarePhotoResponse struct {
	Response PhotoResponse `json:"response"`
}

type Photo struct {
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
}

type Response struct {
	Venues []Venue `json:"venues"`
}

type Venue struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Contact Contact `json:"contact"`
	Location Location `json:"location"`
	Photos []Photo `json:"photos"`
	Categories []Category `json:"categories"`
	Rating float64 `json:"rating"`
	Url string `json:"url"`
}

type Location struct {
	FormattedAddress []string `json:"formattedAddress"`
}

type Contact struct{
	Phone string `json:"phone"`
}

type Category struct {
	Name string `json:"name"`
	Icon CategoryIcon `json:"icon"`
}

type CategoryIcon struct {
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
}

type FourSquareRecommendResponse struct {
	Response GroupResponse `json:"response"`
}

type GroupResponse struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Items []GroupItem `json:"items"`
}

type GroupItem struct {
	Venue Venue `json:"venue"`
}

func (api *FourSquareApi) InitClient() {
	api.Request = resty.R().SetQueryParams(map[string]string{
		"client_id": os.Getenv("API_FOURSQUARE_CLIENT_ID"),
		"client_secret": os.Getenv("API_FOURSQUARE_CLIENT_SECRET"),
		"v" : "20170801",
	}).SetHeader("Accept", "application/json")
}


func (api *FourSquareApi) GetVenue(venueId string) FourSquarePhotoResponse {
	resp, err := api.Request.Get("https://api.foursquare.com/v2/venues/"+venueId+"/photos")

	if err != nil {
		fmt.Println(err)
	}

	var r FourSquarePhotoResponse

	json.Unmarshal(resp.Body(), &r)

	return  r
}

func (api *FourSquareApi) GetVenues(query string) FourSquareResponse {
	resp, err := api.Request.
		SetQueryParams(map[string]string{
		"ll": "59.93,60.00",
		"query" : query,
	}).Get("https://api.foursquare.com/v2/venues/search")

	if err != nil {
		fmt.Println(err)
	}

	var r FourSquareResponse

	json.Unmarshal(resp.Body(), &r)

	return  r
}

func (api *FourSquareApi) GetRecommends(query string) FourSquareRecommendResponse {
	resp, err := api.Request.
		SetQueryParams(map[string]string{
		"sw": "59.843090154492366,29.907188415527344",
		"ne": "59.97425688709357,30.747642517089844",
		"limit": "5",
		"locale": "ru",
		"query" : query,
	}).Get("https://api.foursquare.com/v2/venues/explore")

	if err != nil {
		fmt.Println(err)
	}

	var r FourSquareRecommendResponse

	json.Unmarshal(resp.Body(), &r)

	return  r
}


