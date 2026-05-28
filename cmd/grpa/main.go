package main

import (
	"fmt"
	"log"

	"gitlab.com/victorAFRodrigues/grpa/internal/storage/entities"
	"gitlab.com/victorAFRodrigues/grpa/internal/storage/interfaces"
	"gitlab.com/victorAFRodrigues/grpa/internal/storage/sqlite"
)

func main() {
	db := sqlite.ConnectDatabase()
	defer db.Close()

	var taskRepo interfaces.IRepository[entities.Task]
	taskRepo = sqlite.NewRepository[entities.Task](db, "tasks")

	var configRepo interfaces.IRepository[entities.Config]
	configRepo = sqlite.NewRepository[entities.Config](db, "configs")

	// 1. Testando o repositório de TAREFAS
	tarefas, err := taskRepo.Get()
	if err != nil {
		log.Println("Erro ao buscar tarefas:", err)
	}
	fmt.Printf("Encontradas %d tarefas para processar.\n", len(tarefas))

	// 2. Testando o repositório de CONFIGURAÇÕES
	configuracoes, err := configRepo.Get()
	if err != nil {
		log.Println("Erro ao buscar configurações:", err)
	}
	fmt.Printf("Carregadas %d configurações do sistema.\n", len(configuracoes))
}
