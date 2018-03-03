package main

import (
	"fmt"
	"CookieMonster/modules"
)

func main() {
	ctx, client := modules.ConnectClient()
	file := modules.OpenFile("https://scontent-arn2-1.cdninstagram.com/vp/6e87a7074511a41e6329d463607a129e/5B47771D/t51.2885-15/e35/27893802_258873334651366_2311175026428084224_n.jpg")

	vision := modules.VisionImage{
		Client: client,
		Context: ctx,
		Reader: file,
	}

	fmt.Println(vision.GetLabels())

	//var gisApi modules.GisApi
	//
	//fmt.Println(gisApi.GetItems())
}
