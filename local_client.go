package gobalena

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"
)

type LocalClient interface {
	RestartService(ctx context.Context, serviceName string) error
	StopService(ctx context.Context, serviceName string) error
	StartService(ctx context.Context, serviceName string) error
	ServicesStatus(ctx context.Context) (*Status, error)
	UpdateRelease(ctx context.Context, force bool) error
	RebootSystem(ctx context.Context, force bool) error
	ShutdownSystem(ctx context.Context) error
	ServicesState(ctx context.Context) (*map[string]interface{}, error)
	DeviceState(ctx context.Context) (*DeviceState, error)
	Purge(ctx context.Context) error
	StreamLogs(ctx context.Context, stream chan []byte) error
}

type localClient struct {
	apiKey        string
	supervisorURL string
	supervisorKey string
	appID         string

	httpClient *SturdyClient
}

func NewLocalClient(apiKey, supervisorURL, supervisorKey, appID string) LocalClient {
	return &localClient{
		apiKey:        apiKey,
		supervisorURL: supervisorURL,
		supervisorKey: supervisorKey,
		appID:         appID,

		httpClient: NewSturdyHTTPClient().SetBaseURL(supervisorURL),
	}
}

func (b *localClient) RestartService(ctx context.Context, serviceName string) error {
	err := Unlock(BalenaLockFile)
	if err != nil {
		return fmt.Errorf("error unlocking lockfile before restarting service: %w", err)
	}
	defer func() {
		err = Lock(BalenaLockFile)
		if err != nil {
			fmt.Println("error creating lockfile after restarting service")
		}
	}()

	var data = strings.NewReader(`{"serviceName": "` + serviceName + `"}`)
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(data).
		Post("/v2/applications/" + b.appID + "/restart-service?apikey=" + b.supervisorKey)
	if err != nil {
		return fmt.Errorf("failed performing request to restart service: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("error restarting service: %s", response.Body())
	}

	return nil
}

func (b *localClient) StopService(ctx context.Context, serviceName string) error {
	err := Unlock(BalenaLockFile)
	if err != nil {
		return fmt.Errorf("error unlocking lockfile before stopping service: %w", err)
	}
	defer func() {
		err = Lock(BalenaLockFile)
		if err != nil {
			fmt.Println("error creating lockfile after stopping service")
		}
	}()

	var data = strings.NewReader(`{"serviceName": "` + serviceName + `"}`)
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(data).
		Post("/v2/applications/" + b.appID + "/stop-service?apikey=" + b.supervisorKey)
	if err != nil {
		return fmt.Errorf("failed performing request to stop service: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("error stopping service: %s", response.Body())
	}

	return nil
}

func (b *localClient) StartService(ctx context.Context, serviceName string) error {
	err := Unlock(BalenaLockFile)
	if err != nil {
		return fmt.Errorf("error unlocking lockfile before starting service: %w", err)
	}
	defer func() {
		err = Lock(BalenaLockFile)
		if err != nil {
			fmt.Println("error creating lockfile after starting service")
		}
	}()

	var data = strings.NewReader(`{"serviceName": "` + serviceName + `"}`)
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(data).
		Post("/v2/applications/" + b.appID + "/start-service?apikey=" + b.supervisorKey)
	if err != nil {
		return fmt.Errorf("failed performing request to start service: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("error starting service: %s", response.Body())
	}

	return nil
}

func (b *localClient) ServicesStatus(ctx context.Context) (*Status, error) {
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Status{}).
		Get("/v2/state/status?apikey=" + b.supervisorKey)
	if err != nil {
		return nil, fmt.Errorf("failed performing request for services status: %w", err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting services status: %s", response.Body())
	}

	balenaResult := response.Result().(*Status)
	return balenaResult, nil
}

func (b *localClient) UpdateRelease(ctx context.Context, force bool) error {
	err := Unlock(BalenaLockFile)
	if err != nil {
		return fmt.Errorf("error unlocking lockfile before updating release: %w", err)
	}
	defer func() {
		err = Lock(BalenaLockFile)
		if err != nil {
			fmt.Println("error creating lockfile after updating release")
		}
	}()

	var data = strings.NewReader(`{"force": "` + strconv.FormatBool(force) + `"}`)
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(data).
		Post("/v1/update?apikey=" + b.supervisorKey)
	if err != nil {
		return fmt.Errorf("failed performing request for updating release: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("error updating release: %s", response.Body())
	}

	return nil
}

func (b *localClient) RebootSystem(ctx context.Context, force bool) error {
	err := Unlock(BalenaLockFile)
	if err != nil {
		return fmt.Errorf("error unlocking lockfile before rebooting system: %w", err)
	}

	var data = strings.NewReader(`{"force": "` + strconv.FormatBool(force) + `"}`)
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(data).
		Post("/v1/reboot?apikey=" + b.supervisorKey)
	if err != nil {
		return fmt.Errorf("failed performing request for rebooting system: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("error rebooting system: %s", response.Body())
	}

	return nil
}

func (b *localClient) ShutdownSystem(ctx context.Context) error {
	err := Unlock(BalenaLockFile)
	if err != nil {
		return fmt.Errorf("error unlocking lockfile before shutting down: %w", err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		Post("/v1/shutdown?apikey=" + b.supervisorKey)
	if err != nil {
		return fmt.Errorf("failed performing request for shutting system down: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("error shutting system down: %s", response.Body())
	}

	return nil
}

func (b *localClient) ServicesState(ctx context.Context) (*map[string]interface{}, error) {
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(map[string]interface{}{}).
		Get("/v2/applications/state?apikey=" + b.supervisorKey)
	if err != nil {
		return nil, fmt.Errorf("failed performing request to get services state: %w", err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting services state: %s", response.Body())
	}

	balenaResult := response.Result().(*map[string]interface{})
	return balenaResult, nil
}

func (b *localClient) DeviceState(ctx context.Context) (*DeviceState, error) {
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(DeviceState{}).
		Get("/v1/device?apikey=" + b.supervisorKey)
	if err != nil {
		return nil, fmt.Errorf("failed performing request to get device state: %w", err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting device state: %s", response.Body())
	}

	balenaResult := response.Result().(*DeviceState)
	return balenaResult, nil
}

func (b *localClient) Purge(ctx context.Context) error {
	err := Unlock(BalenaLockFile)
	if err != nil {
		return fmt.Errorf("error unlocking lockfile before shutting down: %w", err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		Post("/v2/applications/" + b.appID + "/purge?apikey=" + b.supervisorKey)
	if err != nil {
		return fmt.Errorf("failed performing request to purge: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("error purging: %s", response.Body())
	}

	return nil
}

func (b *localClient) StreamLogs(ctx context.Context, stream chan []byte) error {
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		SetBody(`{"follow":true,"all":true,"unit":"balena.service","count":40,"format":"json"}`).
		Post("/v2/journal-logs?apikey=" + b.supervisorKey)
	if err != nil {
		return fmt.Errorf("failed performing request to stream logs: %w", err)
	}
	defer response.RawResponse.Body.Close()

	if response.IsError() {
		return fmt.Errorf("error streaming logs: %s", response.Body())
	}

	scanner := bufio.NewScanner(response.RawResponse.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			select {
			case <-ctx.Done():
				return ctx.Err()
			case stream <- scanner.Bytes():
			}
		}
	}
	// Check for errors from the scanner
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from scanner: %w", err)
	}

	return nil
}
