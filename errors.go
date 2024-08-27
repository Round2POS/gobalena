package gobalena

import "errors"

var (
	ErrInvalidBalenaDeviceUUID = errors.New("invalid balena device uuid")

	ErrResourceNotFound  = errors.New("resource not found")
	ErrExpectedOneResult = errors.New("expected one result")
	ErrEnvVarNotFound    = errors.New("env var not found")
)
