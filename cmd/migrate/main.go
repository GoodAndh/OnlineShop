package main

import (
	"log"
	"os"
	"restful/app"

	myCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := app.MysqlDb(myCfg.Config{
		User:   "root",
		Passwd: "r23password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "database_baru",

		AllowNativePasswords: true,
		ParseTime:            true,
		MultiStatements:      true,
	})
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://.", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		err := m.Up()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("berhasil up")
	}
	if cmd == "down" {
		err := m.Down()
		if err != nil && err == migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("berhasil down")
	}
	// migrate database "mysql://root@tcp(localhost:3306)/database_baru1" -path cmd/migrate up
	//migrate create -ext sql -dir [cmd/migrate] [namafile]

}
