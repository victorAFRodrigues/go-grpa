package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type RepositorySQLite[T any] struct {
	db        *sqlx.DB
	tableName string
}

func NewRepository[T any](db *sqlx.DB, tableName string) *RepositorySQLite[T] {
	return &RepositorySQLite[T]{db: db, tableName: tableName}
}

func (r *RepositorySQLite[T]) Get() ([]T, error) {
	var list []T

	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	
	err := r.db.Select(&list, query)
	return list, err
}

func (r *RepositorySQLite[T]) Save(entity T) error {
	query := fmt.Sprintf("INSERT INTO %s (id) VALUES (:id)", r.tableName) 
	_, err := r.db.NamedExec(query, entity)
	return err
}

func (r *RepositorySQLite[T]) GetByID(id string) (T, error) {
	var entity T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", r.tableName)
	err := r.db.Get(&entity, query, id)
	return entity, err
}

func (r *RepositorySQLite[T]) Update(entity T) error {
	query := fmt.Sprintf("UPDATE %s SET /* campos a atualizar */ WHERE id = :id", r.tableName)
	_, err := r.db.NamedExec(query, entity)
	return err
}

func (r *RepositorySQLite[T]) Delete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.tableName)
	_, err := r.db.Exec(query, id)
	return err
}
