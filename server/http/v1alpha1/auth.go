package v1alpha1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
)

const (
	HeaderKeyUsername = "X-Username"
	HeaderKeyUserID   = "X-UserID"
)

func (s *Server) CurrentUser(r *web.Ctx) (uuid.UUID, string, error) {
	cookie, err := r.Request.Cookie("splendor_user")
	if err != nil {
		return nil, "", fmt.Errorf("missing username")
	}
	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, "", fmt.Errorf("invalid cookie")
	}
	var p game.Player
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, "", fmt.Errorf("invalid cookie")
	}
	// username, err := r.HeaderValue(HeaderKeyUsername)
	// if err != nil {
	// 	return nil, "", fmt.Errorf("missing username")
	// }
	// if len(username) == 0 {
	// 	return nil, "", fmt.Errorf("missing username")
	// }

	// userID, err := uuidHeaderValue(r, HeaderKeyUserID)
	// if err != nil {
	// 	return nil, "", fmt.Errorf("missing user id")
	// }

	expected := s.GetUserID(p.Username)
	p.ID = s.GetUserID(p.Username)
	// if !userID.Equal(expected) {
	// 	return nil, "", fmt.Errorf("invalid user id")
	// }

	return expected, p.Username, err
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
