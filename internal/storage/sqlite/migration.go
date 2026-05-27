package sqlite

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" 
	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) {
	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Erro ao criar driver de migração: %v", err)
	}

	caminhoMigrations := "file://internal/storage/migrations"

	m, err := migrate.NewWithDatabaseInstance(caminhoMigrations, "sqlite", driver)
	if err != nil {
		log.Fatalf("Erro ao inicializar o Migrate: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Erro crítico ao aplicar as migrations: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Banco de dados já está na versão mais recente.")
	} else {
		fmt.Println("Tabelas criadas/atualizadas com sucesso via Migrations!")
	}
}