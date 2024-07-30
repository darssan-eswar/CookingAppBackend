package util

import (
	"github.com/joho/godotenv"
)

func LoadEnvFromPath(path string) error {
	return godotenv.Load(path)
}

type Optional[T any] struct {
	Value     T
	isPresent bool
}

func (o Optional[T]) IsPresent() bool {
	return o.isPresent
}

func (o Optional[T]) IsEmpty() bool {
	return !o.isPresent
}

func (o Optional[T]) Get() T {
	return o.Value
}
