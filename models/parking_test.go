package models_test

import (
	"testing"

	"br.com.mlabs/models"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	request := models.ParkingRequest{
		Plate: "NOTVALID",
	}

	assert.Equal(t, models.Validate(request), false)

	request.Plate = "ABC1234"
	assert.Equal(t, models.Validate(request), false)

	request.Plate = "ABC-1234"
	assert.Equal(t, models.Validate(request), true)
}
