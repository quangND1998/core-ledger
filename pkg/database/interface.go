package database

type Getter[T any] interface {
	GetById() T
}

type Updater[T any] interface {
	UpdateById() T
}

type Deleter[T any] interface {
	DeleteById() T
}

type Saver[T any] interface {
	Save(T) error
	Saves(T) error
}
