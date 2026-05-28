package postgres

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

func StartDatabase() *sqlx.DB {
	dir := "./data"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Erro ao criar pasta do banco: %v", err)
	}

	dbPath := filepath.Join(dir, "database.db")

	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco SQLite: %v", err)
	}

	runMigrations(db)

	return db
}