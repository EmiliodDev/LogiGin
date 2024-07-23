package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/EmiliodDev/LogiGin/cmd/api"
    "github.com/EmiliodDev/LogiGin/config"
	"github.com/EmiliodDev/LogiGin/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
    cfg := mysql.Config{
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

    initStorage(db)

    server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)
    if err := server.Run(); err != nil {
        log.Fatal(err)
    } 
}

func initStorage(db *sql.DB) {
    err := db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("DB: Successfully connected")
}
