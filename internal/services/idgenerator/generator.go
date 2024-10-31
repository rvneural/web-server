package idgenerator

import (
	"math/rand"
	"strings"
)

type Generator struct {
	length  int
	letters []rune
}

func New(length int) *Generator {
	return &Generator{
		length:  length,
		letters: []rune("abcdefghijklmnopqrstuvwxyz0123456789-"),
	}
}

func (g *Generator) Generate() string {

	length := rand.Intn(g.length-10) + 10

	id := make([]rune, length)
	for i := range id {
		if i == len(id)-1 || i == 0 {
			id[i] = g.letters[rand.Intn(len(g.letters))]
			for id[i] == '-' {
				id[i] = g.letters[rand.Intn(len(g.letters))]
			}
		}
	}
	return strings.ToLower(string(id))
}
