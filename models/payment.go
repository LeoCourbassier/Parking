package models

import "gorm.io/gorm"

// Payment will hold all the payment information
type Payment struct {
	gorm.Model

	ParkingID uint
	Parking   Parking

	Paid bool
}
