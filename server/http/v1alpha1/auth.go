package v1alpha1

import (
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
)

const (
	HeaderKeyUsername = "X-Username"
	HeaderKeyUserID   = "X-UserID"
)

func (s *Server) CurrentUser(r *web.Ctx) (uuid.UUID, string, error) {
	username, err := r.HeaderValue(HeaderKeyUsername)
	if err != nil {
		return nil, "", fmt.Errorf("missing username")
	}
	if len(username) == 0 {
		return nil, "", fmt.Errorf("missing username")
	}

	// userID, err := uuidHeaderValue(r, HeaderKeyUserID)
	// if err != nil {
	// 	return nil, "", fmt.Errorf("missing user id")
	// }

	expected := s.GetUserID(username)
	// if !userID.Equal(expected) {
	// 	return nil, "", fmt.Errorf("invalid user id")
	// }

	return expected, username, err
}

func (s *Server) GetUserID(username string) uuid.UUID {
	s.usersLock.Lock()
	defer s.usersLock.Unlock()
	id, has := s.Users[username]
	if !has {
		return uuid.Empty()
	}
	return id
}

func (s *Server) GetOrSetUserID(username string) uuid.UUID {
	s.usersLock.Lock()
	defer s.usersLock.Unlock()
	id, has := s.Users[username]
	if !has {
		id = uuid.V4()
		s.Users[username] = id
	}
	return id
}

func uuidHeaderValue(r *web.Ctx, key string) (uuid.UUID, error) {
	id, err := r.HeaderValue(key)
	if err != nil {
		return nil, err
	}
	return uuid.Parse(id)
}
