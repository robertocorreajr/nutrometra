package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	dsn := "pgx5://" + os.Getenv("POSTGRES_USER") + ":" +
		os.Getenv("POSTGRES_PASSWORD") + "@" +
		os.Getenv("POSTGRES_HOST") + ":" +
		os.Getenv("POSTGRES_PORT") + "/" +
		os.Getenv("POSTGRES_DB") + "?sslmode=disable"

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("migrate: new: %v", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrate: up complete")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("migrate down: %v", err)
		}
		log.Println("migrate: down 1 step complete")
	default:
		log.Fatalf("unknown direction: %s", direction)
	}
}
