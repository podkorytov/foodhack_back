package main

import (
	"CookieMonster/modules"
	"github.com/joho/godotenv"
	"log"
	"github.com/gin-gonic/gin"
)

func GetList(c *gin.Context) {
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

	c.JSON(200, gin.H{
		"data": foursquareApi.GetVenues(query),
	})


}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.GET("/get-list", GetList)
	r.Run()
}
