package v1alpha1

import (
	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/apiversions"
)

type Header struct {
	Metadata
	Options HeaderOptions
}

type Metadata struct {
	ID          uuid.UUID
	Origin      uuid.UUID
	Destination uuid.UUID
	Reference   uuid.UUID
	Type        PacketType
	APIVersion  string
}

type HeaderOptions map[string]string

func (h *Header) Values() *HeaderOptions {
	if h == nil {
		d := make(HeaderOptions)
		return &d
	}
	if h.Options == nil {
		h.Options = make(HeaderOptions)
	}
	return &h.Options
}

func (ho *HeaderOptions) Append(other HeaderOptions) {
	if other == nil {
		return
	}
	for k, v := range other {
		(*ho)[k] = v
	}
}

func (ho *HeaderOptions) Add(k, v string) {
	if ho == nil {
		return
	}
	(*ho)[k] = v
}

func (ho HeaderOptions) Value(key string) string {
	if ho == nil {
		return ""
	}
	return ho[key]
}

func NewHeader() Header {
	return Header{
		Metadata: Metadata{
			ID:         uuid.V4(),
			APIVersion: apiversions.V1Alpha1,
		},
		Options: make(HeaderOptions),
	}
}
