package main

import (
	"CookieMonster/modules"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, client := modules.ConnectClient()
	file := modules.OpenFile("http://www.navolne.life/images/201801/840423-1515788124.jpg")

	vision := modules.VisionImage{
		Client: client,
		Context: ctx,
		Reader: file,
	}

	var foursquareApi modules.FourSquareApi

	foursquareApi.InitClient()

	label := vision.GetLabels()[0]

	var googleApi modules.GoogleTranslateApi

	googleApi.InitClient()

	var query string

	if label.LanguageCode != "ru" {
		query = googleApi.Translate(label.Label)
	} else {
		query = label.Label
	}

	fmt.Println(query)

	fmt.Println(foursquareApi.GetVenues(query))
}
