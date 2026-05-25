package storage

import (
	"time"

	"modernc.org/libc/uuid"
)

// Task representa a nossa entidade no sistema
type Task struct {
	ID        UUID
	    `db:"id"`
	UseCase   string    `db:"use_case"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

// TaskRepository define as ações que o resto do robô pode fazer no banco
type TaskRepository interface {
	Salvar(task Task) error
	BuscarTodas() ([]Task, error)
}