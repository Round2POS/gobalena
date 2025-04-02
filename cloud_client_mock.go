package gobalena

import (
	"context"
	"io"
)

var ()

type mockCloudClient struct{}

func NewMockCloudClient() CloudClient {
	return &mockCloudClient{}
}

// CreateDeviceServiceEnvVar implements CloudClient.
func (m *mockCloudClient) CreateDeviceServiceEnvVar(ctx context.Context, balenaDeviceUUID string, name string, serviceInstallID int, value string) error {
	return nil
}

// CreateDeviceEnvVar implements CloudClient.
func (m *mockCloudClient) CreateDeviceEnvVar(ctx context.Context, balenaDeviceUUID string, name string, value string) error {
	return nil
}

// DeleteDeviceEnvVar implements CloudClient.
func (m *mockCloudClient) DeleteDeviceEnvVar(ctx context.Context, balenaDeviceID, envVarID int) error {
	return nil
}

// DeleteDevice implements CloudClient.
func (m *mockCloudClient) DeleteDevice(ctx context.Context, balenaDeviceUUID string) error {
	return nil
}

// DownloadOS implements CloudClient.
func (m *mockCloudClient) DownloadOS(ctx context.Context, writer io.Writer, fleet string, deviceType DeviceType, headerSetter HeaderSetter) (string, error) {
	return "", nil
}

// EnablePublicDeviceURL implements CloudClient.
func (m *mockCloudClient) EnablePublicDeviceURL(ctx context.Context, balenaDeviceUUID string) error {
	return nil
}

// GetDevice implements CloudClient.
func (m *mockCloudClient) GetDevice(ctx context.Context, balenaDeviceUUID string) (*Device, error) {
	return &Device{}, nil
}

// GetDeviceDetails implements CloudClient.
func (m *mockCloudClient) GetDeviceDetails(ctx context.Context, balenaDeviceUUID string) (*Device, error) {
	return &Device{}, nil
}

// GetDeviceEnvVarID implements CloudClient.
func (m *mockCloudClient) GetDeviceEnvVarID(ctx context.Context, balenaDeviceID int, key string) (int, error) {
	return 1, nil
}

// GetDeviceEnvVars implements CloudClient.
func (m *mockCloudClient) GetDeviceEnvVars(ctx context.Context, balenaDeviceUUID string) ([]DeviceEnvVar, error) {
	return []DeviceEnvVar{}, nil
}

// GetDeviceID implements CloudClient.
func (m *mockCloudClient) GetDeviceID(ctx context.Context, balenaDeviceUUID string) (int, error) {
	return 1, nil
}

// GetDeviceServiceEnvVars implements CloudClient.
func (m *mockCloudClient) GetDeviceServiceEnvVars(ctx context.Context, balenaDeviceUUID string) ([]DeviceServiceEnvVar, error) {
	return []DeviceServiceEnvVar{}, nil
}

// GetDevicesDetails implements CloudClient.
func (m *mockCloudClient) GetDevicesDetails(ctx context.Context, balenaDeviceUUIDs []string) ([]Device, error) {
	return []Device{}, nil
}

// GetServiceInstallID implements CloudClient.
func (m *mockCloudClient) GetDeviceServiceInstallIDs(ctx context.Context, balenaDeviceUUID string) ([]DeviceServiceInstall, error) {
	return []DeviceServiceInstall{}, nil
}

// GetFleet implements CloudClient.
func (m *mockCloudClient) GetFleet(ctx context.Context, name string) (*Fleet, error) {
	return &Fleet{}, nil
}

// GetFleetEnvVars implements CloudClient.
func (m *mockCloudClient) GetFleetEnvVars(ctx context.Context, name string) ([]FleetEnvVar, error) {
	return []FleetEnvVar{}, nil
}

// GetFleetReleases implements CloudClient.
func (m *mockCloudClient) GetFleetReleases(ctx context.Context, name string) ([]Release, error) {
	return []Release{}, nil
}

// GetServiceEnvVars implements CloudClient.
func (m *mockCloudClient) GetServiceEnvVars(ctx context.Context, fleetName string) ([]ServiceEnvVar, error) {
	return []ServiceEnvVar{}, nil
}

// DeleteDeviceServiceEnvVar implements CloudClient.
func (m *mockCloudClient) DeleteDeviceServiceEnvVar(ctx context.Context, balenaDeviceID, envVarID int) error {
	return nil
}

// HostLogin implements CloudClient.
func (m *mockCloudClient) HostLogin(token string) error {
	return nil
}

// MoveDeviceToFleet implements CloudClient.
func (m *mockCloudClient) MoveDeviceToFleet(ctx context.Context, balenaDeviceUUID string, fleetName string) error {
	return nil
}

// PinDeviceToRelease implements CloudClient.
func (m *mockCloudClient) PinDeviceToRelease(ctx context.Context, balenaDeviceUUID string, releaseID int) error {
	return nil
}

// RegisterDevice implements CloudClient.
func (m *mockCloudClient) RegisterDevice(ctx context.Context, balenaDeviceUUID string, fleetName string, deviceType DeviceType) error {
	return nil
}

// SetDeviceName implements CloudClient.
func (m *mockCloudClient) SetDeviceName(ctx context.Context, balenaDeviceUUID string, name string) error {
	return nil
}

// UpdateDeviceEnvVar implements CloudClient.
func (m *mockCloudClient) UpdateDeviceEnvVar(ctx context.Context, balenaDeviceID int, envVarID int, value string) error {
	return nil
}

// UpdateDeviceServiceEnvVar implements CloudClient.
func (m *mockCloudClient) UpdateDeviceServiceEnvVar(ctx context.Context, balenaDeviceID int, envVarID int, value string) error {
	return nil
}
