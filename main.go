package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/podkorytov/foodhack_back/modules"
	"net/http"
	"regexp"
)

func GetUrl(c *gin.Context) {
	url := c.Query("url")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Url must be set"})

		return
	}

	match, _ := regexp.MatchString(".(jpeg|jpg|gif|png)", url)

	if match != true {
		instaUrl, err := modules.GetInstaImage(url)
		url = instaUrl

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't parse instagram url"})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func GetRecommends(c *gin.Context) {
	url := c.Query("url")
	ll := c.DefaultQuery("ll", "59.973047,30.340984")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Url must be set"})

		return
	}

	match, _ := regexp.MatchString(".(jpeg|jpg|gif|png)", url)

	if match != true {
		instaUrl, err := modules.GetInstaImage(url)
		url = instaUrl

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't parse instagram url"})

			return
		}
	}

	ctx, client := modules.ConnectClient()
	file := modules.OpenFile(url)

	vision := modules.VisionImage{
		Client:  client,
		Context: ctx,
		Reader:  file,
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

		recommends := foursquareApi.GetRecommends(query, ll)
		categories := foursquareApi.GetCategories()

		for i, group := range recommends.Response.Groups {

			var items []modules.GroupItem

			for _, venue := range group.Items {
				if len(venue.Venue.Categories) == 0 {
					items = append(items, venue)
				}

				if len(venue.Venue.Categories) > 0 {
					if modules.FindCategory(venue.Venue.Categories[0].Id, categories.Categories) == true {
						items = append(items, venue)
					}
				}
			}

			if len(items) > 5 {
				items = items[:5]
			}

			for i, item := range items {
				items[i].Venue.Photos = item.GetPhotos(foursquareApi)
			}

			recommends.Response.Groups[i].Items = items
		}

		c.JSON(200, gin.H{
			"query": query,
			"data":  recommends,
		})

		return
	}

	c.JSON(404, gin.H{
		"message": "Not found",
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if authorization != "SWdlvJ8O9FqwvHpN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})

			c.Abort()
		}

		c.Next()
	}
}

func main() {
	godotenv.Load()

	r := gin.Default()

	r.Use(AuthMiddleware())

	r.GET("/get-recommends", GetRecommends)
	r.GET("/get-url", GetUrl)
	r.Run()
}
