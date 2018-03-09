package modules

import (
	"github.com/utahta/go-openuri"
	"log"

	// Imports the Google Cloud Vision API client package.
	vision "cloud.google.com/go/vision/apiv1"
	"fmt"
	"golang.org/x/net/context"
	vision2 "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"io"
)

type VisionImage struct {
	Reader  io.Reader
	Client  *vision.ImageAnnotatorClient
	Context context.Context
}

func ConnectClient() (context.Context, *vision.ImageAnnotatorClient) {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return ctx, client
}

func OpenFile(url string) io.ReadCloser {
	o, _ := openuri.Open(url)

	return o
}

func (i VisionImage) GetAllLabels() []*vision2.EntityAnnotation {
	image, err := vision.NewImageFromReader(i.Reader)

	if err != nil {
		fmt.Println("Failed to create image: %v", err)
	}

	detect, err := i.Client.DetectLabels(i.Context, image, nil, 10)

	if err != nil {
		fmt.Println("Failed to detect labels: %v", err)
	}

	return detect
}

func (i VisionImage) GetLabels() []*vision2.WebDetection_WebLabel {
	image, err := vision.NewImageFromReader(i.Reader)

	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	detect, err := i.Client.DetectWeb(i.Context, image, nil)

	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	return detect.BestGuessLabels
}
