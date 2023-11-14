package models

import (
	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
)

type Book struct {
	weaver.AutoMarshal
	ID         uuid.UUID
	Author     string
	Title      string
	Descrition string
}
