package sqlite

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) error {

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("falha ao criar driver de migração: %w", err)
	}

	migrationsPath := "file://internal/storage/migrations"

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("falha ao inicializar o Migrate: %w", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Banco de dados já está na versão mais recente.")
			return nil
		}
		return fmt.Errorf("erro ao aplicar migrations: %w", err)
	}

	fmt.Println("Migrations aplicadas com sucesso!")
	return nil
}