package usecases

import (
	"unicode"

	"br.com.mlabs/models"
	"br.com.mlabs/utils"
	"github.com/otiai10/gosseract"
)

// Recognize gets words from an image
func Recognize(path string) (models.ParkingRequest, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(path)
	text, err := client.Text()
	if err != nil {
		return models.ParkingRequest{}, utils.ErrImageRecognition
	}

	cleanedText := ""
	for _, s := range text {
		if unicode.IsLetter(s) || unicode.IsDigit(s) || s == '-' {
			cleanedText += string(s)
		}
	}

	return models.ParkingRequest{
		Plate: cleanedText,
	}, nil
}
