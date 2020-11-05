package usecases_test

import (
	"testing"

	"br.com.mlabs/models"
	"br.com.mlabs/usecases"
	"br.com.mlabs/utils"
	"github.com/stretchr/testify/assert"
)

func TestRecognize(t *testing.T) {
	_, err := usecases.Recognize("notfound")

	assert.Equal(t, err, utils.ErrImageRecognition)

	request, _ := usecases.Recognize("../assets/download.jpg")
	assert.Equal(t, request, models.ParkingRequest{
		Plate: "GTJ-6699",
	})

	request, _ = usecases.Recognize("../assets/download2.jpg")
	assert.Equal(t, request, models.ParkingRequest{
		Plate: "BRA3R52",
	})
}
