package usecases

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"br.com.mlabs/models"
	"br.com.mlabs/storage"
	"br.com.mlabs/utils"
)

// MakeReservation asserts business logic
func MakeReservation(request models.ParkingRequest) (uint, error) {
	if !models.Validate(request) {
		return 0, utils.ErrPlateNotValid
	}

	return storage.ParkingReservation(request)
}

// GetReservations gets all the reservations under a plante
func GetReservations(request models.ParkingRequest) (models.ParkingHistory, error) {
	if !models.Validate(request) {
		return nil, utils.ErrPlateNotValid
	}

	parkingPayments, err := storage.ParkingHistory(request)
	if err != nil {
		return nil, err
	}

	history := models.ParkingHistory{}
	for _, parking := range parkingPayments {
		left := true
		if parking.Checkout == nil {
			left = false
			now := time.Now()
			parking.Checkout = &now
		}

		timeDiff := fmt.Sprintf("%.0f minutes", parking.Checkout.Sub(parking.Checkin).Minutes())

		entry := models.ParkingHistoryEntry{
			ID:   parking.ID,
			Left: left,
			Paid: parking.Paid,
			Time: timeDiff,
		}
		history = append(history, entry)
	}

	return history, nil
}

// Pay sets the payment as true
func Pay(idVar string) error {
	id, err := strconv.ParseUint(idVar, 10, 64)
	if err != nil {
		return utils.ErrIDNotValid
	}

	return storage.Pay(uint(id))
}

// Checkout checks out a parking space
func Checkout(idVar string) error {
	id64, err := strconv.ParseUint(idVar, 10, 64)
	if err != nil {
		return utils.ErrIDNotValid
	}

	id := uint(id64)

	paid, err := storage.IsPaid(id)
	if err != nil {
		if !errors.Is(err, utils.ErrNotFound) {
			return utils.ErrInternalServer
		}
		return err
	}
	if !paid {
		return utils.ErrPayFirst
	}

	return storage.Checkout(id)
}
