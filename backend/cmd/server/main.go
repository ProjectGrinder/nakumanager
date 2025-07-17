package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func runMigrations() {
	m, err := migrate.New(
		"file://db/migrations",
		"sqlite://app.db",
	)
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal("migration up failed:", err)
	}
}

func main() {
	app := fiber.New()

	conn, err := sql.Open("sqlite", "./app.db")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	_, err = conn.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("failed to enable foreign keys:", err)
	}

	runMigrations()

	SetUpRouters(app, conn)

	log.Fatal(app.Listen(":8080"))
}
