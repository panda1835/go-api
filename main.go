package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panda1835/go-api/detect_api"
	"github.com/panda1835/go-api/image_processing"
)

func main() {
	router := gin.Default()
	router.POST("/upload", postImage)
	router.Run("localhost:8080")
}

func postImage(c *gin.Context) {
	img_dest := "upload_image/recent_image"

	// fetch image
	file, _ := c.FormFile("file")

	// save image
	c.SaveUploadedFile(file, img_dest)

	var predictions, err = detect_api.DetectLabels(img_dest)

	// defer os.Remove(file.Filename)

	if err != nil {
		log.Fatal(err)
	} else {
		c.IndentedJSON(http.StatusAccepted, predictions)
	}

	for i := 0; i < len(predictions); i++ {
		crop_file_name := fmt.Sprintf("crop_image/recent_crop_image_%d", i)
		err = image_processing.CropBoundingBox(img_dest, crop_file_name, predictions[i].Coords)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("done")
}
