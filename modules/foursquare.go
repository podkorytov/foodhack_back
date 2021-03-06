package modules

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"io/ioutil"
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
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Contact    Contact    `json:"contact"`
	Location   Location   `json:"location"`
	Photos     []Photo    `json:"photos"`
	Categories []Category `json:"categories"`
	Rating     float64    `json:"rating"`
	Url        string     `json:"url"`
}

type Location struct {
	FormattedAddress []string `json:"formattedAddress"`
	Lat              float64  `json:"lat"`
	Lng              float64  `json:"lng"`
}

type Contact struct {
	Phone string `json:"phone"`
}

type Category struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Icon       CategoryIcon `json:"icon"`
	Categories []Category   `json:"categories"`
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
	Type  string      `json:"type"`
	Name  string      `json:"name"`
	Items []GroupItem `json:"items"`
}

type GroupItem struct {
	Venue Venue `json:"venue"`
}

func (groupItem *GroupItem) GetPhotos(api FourSquareApi) []Photo {
	photo := api.GetVenue(groupItem.Venue.Id)
	photos := photo.Response.Photos.Items

	if len(photos) > 5 {
		photos = photos[:5]
	}

	return photos
}

func (api *FourSquareApi) InitClient() {
	api.Request = resty.R().SetQueryParams(map[string]string{
		"client_id":     os.Getenv("API_FOURSQUARE_CLIENT_ID"),
		"client_secret": os.Getenv("API_FOURSQUARE_CLIENT_SECRET"),
		"v":             "20170801",
	}).SetHeader("Accept", "application/json")
}

func (api *FourSquareApi) GetVenue(venueId string) FourSquarePhotoResponse {
	resp, err := api.Request.Get("https://api.foursquare.com/v2/venues/" + venueId + "/photos")

	if err != nil {
		fmt.Println(err)
	}

	var r FourSquarePhotoResponse

	json.Unmarshal(resp.Body(), &r)

	return r
}

func (api *FourSquareApi) GetVenues(query string) FourSquareResponse {
	resp, err := api.Request.
		SetQueryParams(map[string]string{
			"ll":    "59.93,60.00",
			"query": query,
		}).Get("https://api.foursquare.com/v2/venues/search")

	if err != nil {
		fmt.Println(err)
	}

	var r FourSquareResponse

	json.Unmarshal(resp.Body(), &r)

	return r
}

func (api *FourSquareApi) GetRecommends(query string, ll string) FourSquareRecommendResponse {
	resp, err := api.Request.
		SetQueryParams(map[string]string{
			"ll":             ll,
			"limit":          "100",
			"locale":         "ru",
			"sortByDistance": "1",
			"radius":         "5000",
			"query":          query,
		}).Get("https://api.foursquare.com/v2/venues/explore")

	if err != nil {
		fmt.Println(err)
	}

	var r FourSquareRecommendResponse

	json.Unmarshal(resp.Body(), &r)

	return r
}

func (api *FourSquareApi) GetCategories() Category {
	file, e := ioutil.ReadFile("./synthetic_db/categories.json")

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var categories Category

	err := json.Unmarshal(file, &categories)

	if err != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	return categories
}

func FindCategory(categoryId string, categories []Category) bool {
	for _, category := range categories {
		if category.Id == categoryId {
			return true
		}

		if len(category.Categories) > 0 {
			found := FindCategory(categoryId, category.Categories)

			if found == true {
				return found
			}
		}
	}

	return false
}
