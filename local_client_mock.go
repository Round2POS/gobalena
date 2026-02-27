package gobalena

import "context"

type mockLocalClient struct{}

func NewMockLocalClient() LocalClient {
	return &mockLocalClient{}
}

// DeviceState implements LocalClient.
func (m *mockLocalClient) DeviceState(ctx context.Context) (*DeviceState, error) {
	return &DeviceState{}, nil
}

// Purge implements LocalClient.
func (m *mockLocalClient) Purge(ctx context.Context) error {
	return nil
}

// RebootSystem implements LocalClient.
func (m *mockLocalClient) RebootSystem(ctx context.Context, force bool) error {
	return nil
}

// RestartService implements LocalClient.
func (m *mockLocalClient) RestartService(ctx context.Context, serviceName string) error {
	return nil
}

// ServicesState implements LocalClient.
func (m *mockLocalClient) ServicesState(ctx context.Context) (*map[string]interface{}, error) {
	return &map[string]interface{}{}, nil
}

// ServicesStatus implements LocalClient.
func (m *mockLocalClient) ServicesStatus(ctx context.Context) (*Status, error) {
	return &Status{}, nil
}

// ShutdownSystem implements LocalClient.
func (m *mockLocalClient) ShutdownSystem(ctx context.Context) error {
	return nil
}

// StartService implements LocalClient.
func (m *mockLocalClient) StartService(ctx context.Context, serviceName string) error {
	return nil
}

// StopService implements LocalClient.
func (m *mockLocalClient) StopService(ctx context.Context, serviceName string) error {
	return nil
}

// StreamLogs implements LocalClient.
func (m *mockLocalClient) StreamLogs(ctx context.Context, stream chan []byte) error {
	return nil
}

// UpdateRelease implements LocalClient.
func (m *mockLocalClient) UpdateRelease(ctx context.Context, force bool) error {
	return nil
}
