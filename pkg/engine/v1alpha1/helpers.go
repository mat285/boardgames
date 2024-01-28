package v1alpha1

import "github.com/blend/go-sdk/uuid"

func excludeUUID(id uuid.UUID, exclude ...uuid.UUID) bool {
	for _, e := range exclude {
		if id.Equal(e) {
			return true
		}
	}
	return false
}
