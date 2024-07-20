package utils

import (
	"encoding/json"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func Decode[T any](bytes []byte) (T, error) {
	var decoded T
	err := json.Unmarshal(bytes, decoded)
	if err != nil {
		return decoded, err
	}
	return decoded, nil
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
