package main

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gin-gonic/gin"
)

func init() {
	// Refer to these functions so that goimports is happy before boilerplate is inserted.
	_ = context.Background()
	_ = vision.ImageAnnotatorClient{}
	_ = os.Open
}

type result struct {
	Name   string      `json:"name"`
	Score  float32     `json:"score"`
	Coords [4][2]int16 `json:"coords"`
}

// detectLabels gets labels from the Vision API for an image at the given file path.
func detectLabels(file string) ([]result, error) {

	var prediction []result

	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return prediction, err
	}

	f, err := os.Open(file)
	if err != nil {
		return prediction, err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return prediction, err
	}
	annotations, err := client.LocalizeObjects(ctx, image, nil)
	if err != nil {
		return prediction, err
	}

	if len(annotations) == 0 {
		fmt.Println("No objects found.")
		return prediction, nil
	}

	img_width, img_height := getSize(file)
	for i, annotation := range annotations {
		prediction = append(prediction, result{})
		prediction[i].Name = annotation.Name
		prediction[i].Score = annotation.Score
		for j, v := range annotation.BoundingPoly.NormalizedVertices {
			prediction[i].Coords[j] = [2]int16{int16(v.X * float32(img_width)), int16(v.Y * float32(img_height))}
		}
	}

	return prediction, nil
}

func main() {
	router := gin.Default()
	router.POST("/upload", postImage)
	router.Run("localhost:8080")
}

func getSize(file string) (int16, int16) {
	fileObject, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	image, _, err := image.Decode(fileObject)

	if err != nil {
		log.Fatal(err)
	}

	return int16(image.Bounds().Dx()), int16(image.Bounds().Dy())
}

func postImage(c *gin.Context) {
	img_dest := "upload_image/recent_image"

	// fetch image
	file, _ := c.FormFile("file")

	// save image
	c.SaveUploadedFile(file, img_dest)

	var predictions, err = detectLabels(img_dest)

	// defer os.Remove(file.Filename)

	if err != nil {
		log.Fatal(err)
	} else {
		c.IndentedJSON(http.StatusAccepted, predictions)
	}
}
