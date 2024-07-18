package main

import (
	"log"
	"os"

	"github.com/EmiliodDev/todoAPI/config"
	"github.com/EmiliodDev/todoAPI/db"
	_ "github.com/go-sql-driver/mysql"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    cfg := mysqlDriver.Config{
        User:                   config.Envs.DBUser,
        Passwd:                 config.Envs.DBPassword,
        Addr:                   config.Envs.DBAddress,
        DBName:                 config.Envs.DBName,
        Net:                    "tcp",
        AllowNativePasswords:   true,
        ParseTime:              true,
    }

    db, err := db.NewSQLStorage(cfg)
    if err != nil {
        log.Fatal(err)
    }

    driver, err := mysqlMigrate.WithInstance(db, &mysqlMigrate.Config{})
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://db/migrate/migrations",
        "mysql",
        driver,
    )
    if err != nil {
        log.Fatal(err)
    }

    v, d, _ := m.Version()
    log.Printf("Version: %d, dirty: %v", v, d)

    cmd := os.Args[len(os.Args)-1]
    if cmd == "up" {
        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
    }
    if cmd == "down" {
        if err := m.Down(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
    }
}

