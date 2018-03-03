package main

import (
	"github.com/podkorytov/foodhack_back/modules"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

func GetList(c *gin.Context) {

	url := c.Query("url")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Url must be set"})

		return
	}

	ctx, client := modules.ConnectClient()
	file := modules.OpenFile(url)

	vision := modules.VisionImage{
		Client: client,
		Context: ctx,
		Reader: file,
	}

	var foursquareApi modules.FourSquareApi

	foursquareApi.InitClient()

	labels := vision.GetLabels()

	if len(labels) > 0 {
		label := labels[0]

		var googleApi modules.GoogleTranslateApi

		googleApi.InitClient()

		var query string

		if label.LanguageCode != "ru" {
			query = googleApi.Translate(label.Label)
		} else {
			query = label.Label
		}

		c.JSON(200, gin.H{
			"query" : query,
			"data": foursquareApi.GetVenues(query),
		})

		return
	}

	c.JSON(404, gin.H{
		"message" : "Not found",
	})
}

func main() {
	godotenv.Load()

	r := gin.Default()

	r.GET("/get-list", GetList)
	r.Run()
}
