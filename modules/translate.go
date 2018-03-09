package modules

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"os"
)

type GoogleTranslateApi struct {
	Request *resty.Request
}

type Translation struct {
	TranslatedText string `json:"translatedText"`
}

type Data struct {
	Translations []Translation `json:"translations"`
}

type GoogleTranslateResponse struct {
	Data Data `json:"data"`
}

func getGoogleToken() string {
	resp, err := resty.R().Get(os.Getenv("API_GOOGLE_TOKEN"))

	if err != nil {
		fmt.Println(err)
	}

	return resp.String()
}

func (api *GoogleTranslateApi) InitClient() {
	api.Request = resty.R().SetQueryParams(map[string]string{
		"source": "en",
		"target": "ru",
		"format": "text",
	})
}

func (api *GoogleTranslateApi) Translate(query string) string {
	fmt.Println(query)

	resp, err := api.Request.
		SetQueryParams(map[string]string{
			"q": query,
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken(getGoogleToken()).
		Get("https://translation.googleapis.com/language/translate/v2")

	if err != nil {
		fmt.Println(err)
	}

	var response GoogleTranslateResponse

	json.Unmarshal(resp.Body(), &response)

	if len(response.Data.Translations) > 0 {
		return response.Data.Translations[0].TranslatedText
	}

	return "empty string"
}
