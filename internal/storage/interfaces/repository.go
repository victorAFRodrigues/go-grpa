package interfaces

type IRepository[T any] interface {
	Get() ([]T, error)
	GetByID(id string) (T, error)
	Save(entity T) error
	Update(entity T) error
	Delete(id string) error
}