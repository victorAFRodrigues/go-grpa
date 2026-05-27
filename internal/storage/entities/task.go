package entities

import (
	"time"
)

type TaskStatus = string

const (
	Running  TaskStatus = "E"
	Finished TaskStatus = "F"
	Failed   TaskStatus = "D"
)

type Task struct {
	ID           string     `db:"id"`
	UseCase      string     `db:"use_case"`
	Data         string     `db:"data"`
	ErrorMessage string     `db:"error_message"`
	Status       TaskStatus `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
