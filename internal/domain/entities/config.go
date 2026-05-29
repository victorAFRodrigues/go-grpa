package entities

import (
	"time"
)

type Config struct {
	Id             string    `db:"id"`
	GrpaName       string    `db:"grpa_name"`
	GrpaVersion    string    `db:"grpa_version"`
	WorkerTimeout  int       `db:"worker_timeout"`
	SystemName     string    `db:"system_name"`
	SystemUrl      string    `db:"system_url"`
	SystemUsername string    `db:"system_username"`
	SystemPassword string    `db:"system_password"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
