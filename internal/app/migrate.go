package app

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/viper"
	"log"
	"time"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate() {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.dbname"),
	)
	fmt.Printf("url: %s\n", databaseURL)
	var (
		attempts = 5
		timeout  = time.Second
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		log.Printf("migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(timeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrate: no change")
		return
	}

	log.Printf("migrate: up success")
}
