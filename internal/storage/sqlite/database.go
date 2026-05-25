package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func StartDatabase(db *sqlx.DB) *sqlx.DB {
	database := db.Connect()

	fmt.Println("Banco de dados inicializado e atualizado com sucesso!")

	return database
}