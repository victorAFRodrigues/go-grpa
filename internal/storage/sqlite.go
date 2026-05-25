package storage

import (
	"fmt"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func Connect() *sqlx.DB {
	if err := os.MkdirAll("../data", os.ModePerm); err != nil {
		log.Fatalf("Erro ao criar pasta data: %v", err)
	}

	db, err := sqlx.Connect("sqlite", "../data/database.db")
	if err != nil {
		log.Fatalf("Erro ao abrir banco: %v", err)
	}

	fmt.Println("Checando atualizações no banco de dados...")

	if err := RunMigrations(db); err != nil {
		log.Fatalf("Falha crítica no banco de dados: %v", err)
	}

	return db
}

