/*
 * @Version: 1.0
 * @Date: 2023-02-20 22:04:59
 * @LastEditTime: 2023-03-19 21:07:17
 */
package model

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB = Init()

func Init() *gorm.DB {
	dsn := "user:password@tcp(127.0.0.1:3306)/my_todo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Migrate the schema
	//db.AutoMigrate(&model.Todo{})
	return db
}
