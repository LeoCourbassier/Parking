package storage

import (
	"fmt"
	"os"
	"strings"
	"time"

	"br.com.mlabs/models"
	"br.com.mlabs/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Connect connects to psql database
func Connect() {
	user, ok := os.LookupEnv("DATABASE_USER")
	if !ok {
		logrus.Panic("DATABASE_USER not found")
	}

	passwd, ok := os.LookupEnv("DATABASE_PASS")
	if !ok {
		logrus.Panic("DATABASE_USER not found")
	}

	host, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		logrus.Panic("DATABASE_USER not found")
	}

	port, ok := os.LookupEnv("DATABASE_PORT")
	if !ok {
		logrus.Panic("DATABASE_USER not found")
	}

	schema, ok := os.LookupEnv("DATABASE_SCHEMA")
	if !ok {
		logrus.Panic("DATABASE_USER not found")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, passwd, host, port, schema)
	logrus.Debugf("Connecting to database with: %s", dsn)

	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		logrus.Panic(err.Error())
	}

	migrate()
}

// ConnectTest connects to the test database
func ConnectTest() {
	if db != nil {
		return
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "postgres", "postgres", "localhost", "5432", "mlabs_test")
	logrus.Debugf("Connecting to database with: %s", dsn)

	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		logrus.Panic(err.Error())
	}

	db = db.Debug()
	migrate()
}

// StartTest gets an instance of db
func StartTest() *gorm.DB {
	return db
}

func migrate() {
	db.AutoMigrate(&models.Parking{})
	db.AutoMigrate(&models.Payment{})
}

// ParkingReservation creates a new record on the database
func ParkingReservation(request models.ParkingRequest) (uint, error) {
	parking := models.Parking{
		Plate:   request.Plate,
		Checkin: time.Now(),
	}
	err := db.Create(&parking).Error
	if err != nil {
		logrus.Warn(err.Error())
		return 0, utils.ErrInternalServer
	}

	return parking.ID, nil
}

// ParkingHistory gets all reservation entries
func ParkingHistory(request models.ParkingRequest) ([]models.ParkingPayments, error) {
	rows, err := db.Model(&models.Parking{}).Joins("LEFT JOIN payments ON parkings.id = payments.parking_id").Select("parkings.*, payments.paid").Where("parkings.plate = ?", request.Plate).Rows()
	defer rows.Close()
	if err != nil {
		logrus.Warn(err.Error())
		return nil, utils.ErrInternalServer
	}

	var parkingWithPayments []models.ParkingPayments
	for rows.Next() {
		var parking models.ParkingPayments
		db.ScanRows(rows, &parking)

		parkingWithPayments = append(parkingWithPayments, parking)
	}

	if parkingWithPayments == nil {
		return nil, utils.ErrNotFound
	}

	return parkingWithPayments, nil
}

// Pay sets the payment in the database
func Pay(id uint) error {
	paid, err := IsPaid(id)

	if err != nil {
		return err
	}
	if paid {
		return utils.ErrAlreadyPaid
	}

	pay := models.Payment{
		Paid:      true,
		ParkingID: id,
	}
	err = db.Create(&pay).Error
	if err != nil {
		logrus.Warn(err.Error())
		return utils.ErrInternalServer
	}

	return nil
}

// IsPaid returns true if a parking space has been paid
func IsPaid(id uint) (bool, error) {
	var res models.ParkingPayments
	tx := db.Joins("LEFT JOIN payments ON payments.parking_id = parkings.id").Select("payments.paid").Where("parkings.id = ?", id).First(&models.Parking{})
	if tx.Error != nil {
		if strings.Contains(tx.Error.Error(), "record not found") {
			return false, utils.ErrNotFound
		}
		return false, utils.ErrInternalServer
	}

	tx.Scan(&res)

	return res.Paid, nil
}

// Checkout checks out a parking space
func Checkout(id uint) error {
	ok, err := HaveCheckedOut(id)
	if err != nil {
		return err
	}
	if ok {
		return utils.ErrAlreadyCheckedOut
	}

	err = db.Model(&models.Parking{}).Where("id = ?", id).Update("checkout", time.Now()).Error
	if err != nil {
		logrus.Warn(err.Error())
		return utils.ErrInternalServer
	}

	return nil
}

// HaveCheckedOut returns true if a parking space has been checked out
func HaveCheckedOut(id uint) (bool, error) {
	tx := db.Where("id = ?", id).Select("checkout").First(&models.Parking{})
	if tx.Error != nil {
		logrus.Warn(tx.Error.Error())
		if strings.Contains(tx.Error.Error(), "record not found") {
			return false, utils.ErrNotFound
		}
		return false, utils.ErrInternalServer
	}

	var parking models.Parking
	tx.Scan(&parking)

	return parking.Checkout != nil, nil
}
