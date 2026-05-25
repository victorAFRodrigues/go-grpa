package main

import (
	"gitlab.com/victorAFRodrigues/grpa/internal/storage"
)

func main() {
	db := storage.StartDatabase()
	defer db.Close()

	result, err := db.Exec("select * from tasks")
	if err != nil {
		println("Erro ao executar consulta:", err.Error())
		return
	}
	
	println("result:", result)
}