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

func RandomBalenaUUID() string {
	return FormatBalenaUUID(uuid.New().String())
}

func FormatBalenaUUID(u string) string {
	return strings.ReplaceAll(u, "-", "")
}

func ParseDeviceType(s string) DeviceType {
	return DeviceType(s)
}
