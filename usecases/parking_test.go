package usecases_test

import (
	"fmt"
	"testing"
	"time"

	"br.com.mlabs/models"
	"br.com.mlabs/storage"
	"br.com.mlabs/usecases"
	"br.com.mlabs/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var tx *gorm.DB

func TestMain(t *testing.M) {
	storage.ConnectTest()
	tx = storage.StartTest()

	t.Run()

	tx.Exec("TRUNCATE parkings CASCADE;")
	tx.Exec("ALTER SEQUENCE parkings_id_seq RESTART WITH 1")
	tx.Exec("ALTER SEQUENCE payments_id_seq RESTART WITH 1")
}

func TestPay(t *testing.T) {
	assert.Equal(t, usecases.Pay("notvalid"), utils.ErrIDNotValid)
	assert.Equal(t, usecases.Pay("-1"), utils.ErrIDNotValid)
	assert.Equal(t, usecases.Pay("1"), utils.ErrNotFound)

	// Test for happy path
	parking := models.Parking{
		Plate:   "ABC-1234",
		Checkin: time.Now(),
	}
	tx.Create(&parking)

	assert.Equal(t, usecases.Pay(fmt.Sprint(parking.ID)), nil)

	assert.Equal(t, usecases.Pay(fmt.Sprint(parking.ID)), utils.ErrAlreadyPaid)
}

func TestCheckout(t *testing.T) {
	assert.Equal(t, usecases.Checkout("notvalid"), utils.ErrIDNotValid)
	assert.Equal(t, usecases.Checkout("-1"), utils.ErrIDNotValid)
	assert.Equal(t, usecases.Checkout("1000"), utils.ErrNotFound)

	// Test for happy path
	parking := models.Parking{
		Plate:   "ABC-1234",
		Checkin: time.Now(),
	}
	tx.Create(&parking)

	assert.Equal(t, usecases.Checkout(fmt.Sprint(parking.ID)), utils.ErrPayFirst)

	payment := models.Payment{
		Paid:      true,
		ParkingID: parking.ID,
	}

	tx.Create(&payment)

	assert.Equal(t, usecases.Checkout(fmt.Sprint(parking.ID)), nil)

	assert.Equal(t, usecases.Checkout(fmt.Sprint(parking.ID)), utils.ErrAlreadyCheckedOut)
}

func TestMakeReservation(t *testing.T) {
	parking := models.ParkingRequest{
		Plate: "ABC-1234",
	}

	id, err := usecases.MakeReservation(parking)
	assert.Equal(t, err, nil)
	assert.Greater(t, id, uint(0))

	parking.Plate = "ab"
	id, err = usecases.MakeReservation(parking)
	assert.Equal(t, err, utils.ErrPlateNotValid)
	assert.Equal(t, id, uint(0))
}

func TestGetReservations(t *testing.T) {
	request := models.ParkingRequest{
		Plate: "ABC-12555",
	}

	history, err := usecases.GetReservations(request)
	assert.Equal(t, err, utils.ErrPlateNotValid)
	assert.Equal(t, history, models.ParkingHistory(nil))

	request.Plate = "ABD-9999"
	history, err = usecases.GetReservations(request)
	assert.Equal(t, err, utils.ErrNotFound)
	assert.Equal(t, history, models.ParkingHistory(nil))

	// Happy path, history with no payment
	parking := models.Parking{
		Plate:   request.Plate,
		Checkin: time.Now(),
	}
	tx.Create(&parking)

	history, err = usecases.GetReservations(request)
	assert.Equal(t, history, models.ParkingHistory{
		{
			ID:   parking.ID,
			Time: "0 minutes",
			Paid: false,
			Left: false,
		},
	})
	assert.Equal(t, err, nil)

	// Happy path, history with paid = true
	pay := models.Payment{
		Paid:      true,
		ParkingID: parking.ID,
	}
	tx.Create(&pay)

	history, err = usecases.GetReservations(request)
	assert.Equal(t, history, models.ParkingHistory{
		{
			ID:   parking.ID,
			Time: "0 minutes",
			Paid: true,
			Left: false,
		},
	})
	assert.Equal(t, err, nil)

	// Happy path, history with two entries
	parking2 := models.Parking{
		Plate:   request.Plate,
		Checkin: time.Now(),
	}
	tx.Create(&parking2)

	history, err = usecases.GetReservations(request)
	assert.Equal(t, history, models.ParkingHistory{
		{
			ID:   parking.ID,
			Time: "0 minutes",
			Paid: true,
			Left: false,
		},
		{
			ID:   parking2.ID,
			Time: "0 minutes",
			Paid: false,
			Left: false,
		},
	})
	assert.Equal(t, err, nil)
}
