package models

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Parking is the model that will hold parking information
type Parking struct {
	gorm.Model

	Plate    string    `gorm:"not null;varchar(8)"`
	Checkin  time.Time `sql:"DEFAULT:current_timestamp"`
	Checkout *time.Time
}

// ParkingRequest will hold the parking reservation
type ParkingRequest struct {
	Plate string `json:"plate" validate:"plate"`
}

// ParkingHistoryEntry is a parking history entry
type ParkingHistoryEntry struct {
	ID   uint   `json:"id"`
	Time string `json:"time"`
	Paid bool   `json:"paid"`
	Left bool   `json:"left"`
}

// ParkingPayments is used to begin history assemble
type ParkingPayments struct {
	ID       uint
	Paid     bool
	Checkin  time.Time
	Checkout *time.Time
	Plate    string
}

// ParkingHistory holds all entries for a plate
type ParkingHistory []ParkingHistoryEntry

// Validate returns true if the request is valid
func Validate(model interface{}) bool {
	v := validator.New()
	v.RegisterValidation("plate", plateValidation, false)
	err := v.Struct(model)
	if err != nil {
		valError := err.(validator.ValidationErrors)
		logrus.Warnf("Validation failed for: %v %s", model, valError.Error())
		return false
	}

	return true
}

func plateValidation(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	old, _ := regexp.MatchString(`^[A-Z]{3}-[0-9]{4}\z`, str)
	mercosul, _ := regexp.MatchString(`^[A-Z]{3}[0-9][A-Z][0-9]{2}\z`, str)

	return old || mercosul
}
