package main

import (
	"ddd2/app"
	"ddd2/cmd/api"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

func main() {
	var validate *validator.Validate = validator.New()

	db, err := app.NewMysqlDb(mysql.Config{
		User:                 "root",
		Passwd:               "r23password",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "glg_restful",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewApiServer(":3000",db,validate)
	if err:=server.Run();err!=nil{
		log.Fatal(err)
	}
}
