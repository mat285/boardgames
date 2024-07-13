package v1alpha1

import (
	"context"
	"sync"

	"github.com/blend/go-sdk/uuid"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

var _ ClientInfo = new(MultiConn)

type MultiConn struct {
	sync.Mutex
	ID       uuid.UUID
	Username string
	Pool     map[ClientInfo]bool
}

func NewMulti(id uuid.UUID, username string) *MultiConn {
	return &MultiConn{
		ID:       id,
		Username: username,
		Pool:     make(map[ClientInfo]bool),
	}
}

func (m *MultiConn) GetID() uuid.UUID {
	return m.ID
}

func (m *MultiConn) GetUsername() string {
	return m.Username
}

func (m *MultiConn) Add(ctx context.Context, c ClientInfo) {
	m.Lock()
	defer m.Unlock()
	m.Pool[c] = true
}

func (m *MultiConn) Delete(ctx context.Context, c ClientInfo) {
	m.Lock()
	defer m.Unlock()
	delete(m.Pool, c)
}

func (m *MultiConn) GetPool() []ClientInfo {
	m.Lock()
	defer m.Unlock()
	ret := make([]ClientInfo, 0, len(m.Pool))
	for c := range m.Pool {
		ret = append(ret, c)
	}
	return ret
}

func (m *MultiConn) Send(ctx context.Context, p wire.Packet) error {
	defer func() {
		recover()
	}()

	for _, c := range m.GetPool() {
		if c == nil {
			continue
		}
		c.Send(ctx, p)
	}
	return nil
}
