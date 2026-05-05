package store

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const idAlphabet = "abcdefghijkmnpqrstuvwxyz23456789"

func NewID() string {
	id, err := gonanoid.Generate(idAlphabet, 8)
	if err != nil {
		panic(fmt.Sprintf("nanoid: %v", err))
	}
	return id
}
