package core

import "github.com/google/uuid"

type IDGenerator interface {
	GenerateID() string
}

func NewIDGenerator() IDGenerator {
	return idGenerator{}
}

type idGenerator struct{}

func (g idGenerator) GenerateID() string {
	return uuid.New().String()
}
