package image_processing

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
)

func CropBoundingBox(img_src, img_dest string, coord [4][2]int16) error {

	fileObject, err := os.Open(img_src)

	if err != nil {
		log.Fatal(err)
	}

	defer fileObject.Close()

	dest_img, _, err := image.Decode(fileObject)

	if err != nil {
		log.Fatal("unable to read src image")
	}

	// create a blank canvas for each bounding box
	bounding_canvas := image.NewRGBA(image.Rect(0, 0, int(coord[1][0]-coord[0][0]), int(coord[2][1]-coord[1][1])))

	// draw the bounding section of the image on the just created canvas
	draw.Draw(bounding_canvas, image.Rect(0, 0, int(coord[1][0]-coord[0][0]), int(coord[2][1]-coord[1][1])), dest_img, image.Point{int(coord[0][0]), int(coord[0][1])}, draw.Src)

	file, err := os.Create(img_dest)
	if err != nil {
		log.Fatal(err)
	}

	png.Encode(file, bounding_canvas)

	return err
}
