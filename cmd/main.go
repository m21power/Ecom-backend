package main

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/m21power/Ecom/cmd/api"
	"github.com/m21power/Ecom/config"
	"github.com/m21power/Ecom/db"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		DBName:               config.Envs.DBName,
		Addr:                 config.Envs.DBAddress,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		fmt.Println("Error connecting database")
		return
	}
	server := api.NewAPIServer(":8080", db)
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
