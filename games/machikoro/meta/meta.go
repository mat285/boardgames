package meta

import (
	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

const (
	Name = "machikoro"
)

var (
	ID = uuid.V4()
)

type Object struct {
}

func (o Object) Meta() v1alpha1.Meta {
	return Meta{}
}

type Meta struct {
}

func (m Meta) ID() uuid.UUID {
	return ID
}

func (m Meta) Name() string {
	return Name
}
