package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SensorReading struct {
	Timestamp   int64   `json:"timestamp"`
	Device      string  `json:"device"`
	Adapter     string  `json:"adapter"`
	Measurement string  `json:"measurement"`
	Value       float64 `json:"value"`
}

type Database struct {
	gdb *gorm.DB
}

func (db *Database) Write(s SensorReading) error {
	return db.gdb.Create(&s).Error
}

func (db *Database) Get(from int64, to int64) ([]SensorReading, error) {
	var sensorReadings []SensorReading
	err := db.gdb.
		Where("timestamp between ? and ?", from, to).
		Find(&sensorReadings, "").Error
	if err != nil {
		return nil, err
	}

	return sensorReadings, nil
}

func (db *Database) ClearExpiredRecords() error {
	config := GetConfig()
	expiry := time.Now().UnixMilli() - int64(config.Retention)*1000
	return db.gdb.Exec("delete from sensor_readings where timestamp < ?", expiry).Error
}

var db *Database

func initDB() {
	config := GetConfig()
	gdb, err := gorm.Open(sqlite.Open(config.DatabasePath), &gorm.Config{})
	if err != nil {
		log.Fatalf("error opening database: %s", err.Error())
	}

	var SensorReading SensorReading
	err = gdb.AutoMigrate(SensorReading)
	if err != nil {
		log.Fatalf("error migrating sensor reading: %s", err.Error())
	}

	db = &Database{gdb: gdb}
}

func GetDatabase() *Database {
	if db == nil {
		initDB()
	}

	return db
}
