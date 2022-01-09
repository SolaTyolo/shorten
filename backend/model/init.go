package model

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	Self   *gorm.DB
	Docker *gorm.DB
}

var DB *Database

func (db *Database) Init() {
	DB = &Database{
		Self:   GetSelfDB(),
		Docker: GetDockerDB(),
	}
}

func GetSelfDB() *gorm.DB {
	return InitSelfDB()
}

func GetDockerDB() *gorm.DB {
	return InitDockerDB()
}

func InitSelfDB() *gorm.DB {
	return openDB(
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"),
	)
}

func InitDockerDB() *gorm.DB {
	return openDB(viper.GetString("docker_db.username"),
		viper.GetString("docker_db.password"),
		viper.GetString("docker_db.addr"),
		viper.GetString("docker_db.name"))
}

func openDB(username, pwd, addr, database string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		username,
		pwd,
		addr,
		database,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// fmt.Errorf(err)
	}

	setupDB(db)
	return db
}

func setupDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(2000)
	// sqlDB.SetConnMaxIdleTime(time.Hour)
}

func (db *Database) Close() {
	d, _ := DB.Self.DB()
	d.Close()
	d, _ = DB.Docker.DB()
	d.Close()
}
