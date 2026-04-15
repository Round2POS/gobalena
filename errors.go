package gobalena

import "errors"

var (
	ErrInvalidBalenaDeviceUUID = errors.New("invalid balena device uuid")
	ErrInvalidMerchantUUID     = errors.New("invalid merchant uuid")
	ErrResourceNotFound  = errors.New("resource not found")
	ErrExpectedOneResult = errors.New("expected one result")
	ErrEnvVarNotFound    = errors.New("env var not found")
	ErrInvalidReleaseID  = errors.New("invalid release ID: must be greater than 0")
)
