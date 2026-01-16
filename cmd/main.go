package main

import (
	"database/sql"
	"log"
	"vietio/config"
	"vietio/internal/app"

	_ "github.com/jackc/pgx/v5/stdlib"
)



func main() {
    config := config.Load()

    dbConn, err := sql.Open("pgx", config.Db.Dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer dbConn.Close()

    if err := dbConn.Ping(); err != nil {
        log.Fatal("db ping failed:", err)
    }

    if config.SeedFlag {
        app.RunSeed(dbConn, config)
    } else {
        app.RunMigrations(dbConn)
    }

    app.RunHttpServer(dbConn, config)
}
