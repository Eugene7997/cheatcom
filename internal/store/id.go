package store

import gonanoid "github.com/matoous/go-nanoid/v2"

const idAlphabet = "abcdefghijkmnpqrstuvwxyz23456789"

func NewID() string {
	id, _ := gonanoid.Generate(idAlphabet, 8)
	return id
}
