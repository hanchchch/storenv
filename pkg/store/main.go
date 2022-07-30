package store

type Storage interface {
	Store([]string) error
	Load([]string) error
	Name() string
}
