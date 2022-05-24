package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panda1835/go-api/api"
	"github.com/panda1835/go-api/image_processing"
)

func main() {
	router := gin.Default()
	router.POST("/upload", postImage)
	router.Run("localhost:8080")
}

type post_response struct {
	Name  string  `json:"name"`
	Score float32 `json:"score"`
	URL   string  `json:"crop_url"`
}

func postImage(c *gin.Context) {
	img_dest := "upload_image/recent_image"

	// fetch image
	file, _ := c.FormFile("file")

	// save image
	c.SaveUploadedFile(file, img_dest)

	var predictions, err = api.DetectLabels(img_dest)

	var response []post_response
	for i := 0; i < len(predictions); i++ {
		response = append(response, post_response{})
		// assign prediction name and score
		response[i].Name = predictions[i].Name
		response[i].Score = predictions[i].Score

		// crop image
		crop_file_name := fmt.Sprintf("crop_image/recent_crop_image_%d", i)
		err = image_processing.CropBoundingBox(img_dest, crop_file_name, predictions[i].Coords)
		if err != nil {
			log.Fatal(err)
		}

		// upload to google cloud bucket
		img_url, err := api.UploadFile("go-api", crop_file_name)
		if err != nil {
			log.Fatal(err)
		}

		response[i].URL = img_url
	}

	if err != nil {
		log.Fatal(err)
	} else {
		c.IndentedJSON(http.StatusAccepted, response)
	}

	fmt.Println("done")
}
