package gobalena

import (
	"strings"

	"github.com/google/uuid"
)

func IsValidBalenaDeviceUUID(u string) bool {
	_, err := uuid.Parse(u)
	if err != nil {
		return false
	}

	return !strings.Contains(u, "-")
}
