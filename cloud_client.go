package gobalena

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

const (
	DeviceTypeGeneric   DeviceType = "genericx86-64-ext"
	DeviceTypeSurfaceGo DeviceType = "surface-go"
	DeviceTypeSurface6  DeviceType = "surface-6"

	DeviceQuerySelector           = "$select=id,uuid,ip_address,mac_address,public_address,device_name,os_version,os_variant,supervisor_version,is_online,last_connectivity_event,is_web_accessible,latitude,longitude,location,created_at,overall_status"
	DeviceDetailsQuerySelector    = "$expand=is_running__release($expand=is_created_by__user($select=id,username,created_at),release_tag($select=tag_key,value,id)),should_be_running__release($expand=is_created_by__user($select=id,username,created_at),release_tag($select=tag_key,value,id)),belongs_to__application($select=id,app_name)"
	FleetReleasesQuerySelect      = "$expand=is_created_by__user($select=id,username,created_at),is_running_on__device/$count,release_tag($select=tag_key,value,id)"
	OrderByCreatedAtQuerySelector = "$orderby=created_at%20desc"
)

type CloudClient struct {
	httpClient *SturdyClient
}

func NewCloudClient(apiKey, endpoint string) CloudClient {
	return CloudClient{
		httpClient: NewSturdyHTTPClient().
			SetBaseURL(endpoint).
			SetHeader("Authorization", "Bearer "+apiKey),
	}
}

func (b *CloudClient) GetDevice(
	ctx context.Context,
	balenaDeviceUUID string,
) (*Device, error) {
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[Device]{}).
		Get("/v6/device(uuid='" + balenaDeviceUUID + "')?$select=" + DeviceQuerySelector)
	if err != nil {
		return nil, fmt.Errorf("failed performing request to get device(%s) details: %w", balenaDeviceUUID, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting device(%s) details: %s", balenaDeviceUUID, response.Body())
	}

	balenaResult := response.Result().(*Response[Device])
	if len(balenaResult.D) == 0 {
		return nil, ErrResourceNotFound
	}

	if len(balenaResult.D) > 1 {
		return nil, ErrExpectedOneResult
	}

	return &balenaResult.D[0], nil
}

func (b *CloudClient) GetDeviceEnvVarID(
	ctx context.Context,
	balenaDeviceID int,
	key string,
) (int, error) {
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[DeviceEnvVar]{}).
		Get("/v6/device_environment_variable?\\$filter=device%20eq%20" + strconv.Itoa(balenaDeviceID))
	if err != nil {
		return 0, fmt.Errorf("failed performing request to get device(%d) env var(%s): %w", balenaDeviceID, key, err)
	}

	if response.IsError() {
		return 0, fmt.Errorf("error setting device(%d) env var(%s): %s", balenaDeviceID, key, response.Body())
	}

	balenaResult := response.Result().(*Response[DeviceEnvVar])
	for _, envVar := range balenaResult.D {
		if envVar.Name == key {
			return envVar.ID, nil
		}
	}

	return 0, ErrEnvVarNotFound
}

func (b *CloudClient) UpdateDeviceEnvVar(
	ctx context.Context,
	balenaDeviceID, envVarID int,
	value string,
) error {
	var data = strings.NewReader(`{"value": "` + value + `"}`)
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(data).
		Patch("/v6/device_environment_variable(" + strconv.Itoa(envVarID) + ")?$filter=device%20eq%20" + strconv.Itoa(balenaDeviceID))
	if err != nil {
		return fmt.Errorf("failed performing request to update device(%d) env var(%d): %w", balenaDeviceID, envVarID, err)
	}

	if response.IsError() {
		return fmt.Errorf("error updating device(%d) env var(%d): %s", balenaDeviceID, envVarID, response.Body())
	}

	return nil
}

func (b *CloudClient) GetDeviceDetails(
	ctx context.Context,
	balenaDeviceUUID string,
) (*Device, error) {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return nil, ErrInvalidBalenaDeviceUUID
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[Device]{}).
		Get("/v6/device(uuid='" + balenaDeviceUUID + "')?" + DeviceDetailsQuerySelector + "&" + DeviceQuerySelector)
	if err != nil {
		return nil, fmt.Errorf("failed performing request to get device(%s) details: %w", balenaDeviceUUID, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting device(%s) details: %s", balenaDeviceUUID, response.Body())
	}

	balenaResult := response.Result().(*Response[Device])
	if len(balenaResult.D) == 0 {
		return nil, ErrResourceNotFound
	}

	if len(balenaResult.D) > 1 {
		return nil, ErrExpectedOneResult
	}

	return &balenaResult.D[0], nil
}

func (b *CloudClient) GetDevicesDetails(
	ctx context.Context,
	balenaDeviceUUIDs []string,
) ([]Device, error) {
	filter := ""
	for i, uuid := range balenaDeviceUUIDs {
		if !IsValidBalenaDeviceUUID(uuid) {
			return nil, ErrInvalidBalenaDeviceUUID
		}

		filter += "'" + uuid + "'"
		if i < len(balenaDeviceUUIDs)-1 {
			filter += ","
		}
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[Device]{}).
		Get("/v6/device?$filter=uuid in(" + filter + ")&" + DeviceDetailsQuerySelector + "&" + DeviceQuerySelector)
	if err != nil {
		return nil, fmt.Errorf("failed performing request to get devices(%s) details: %w", balenaDeviceUUIDs, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting devices(%s) details: %s", balenaDeviceUUIDs, response.Body())
	}

	balenaResult := response.Result().(*Response[Device])
	if len(balenaResult.D) == 0 {
		return nil, ErrResourceNotFound
	}

	return balenaResult.D, nil
}

func (b *CloudClient) GetDeviceID(
	ctx context.Context,
	balenaDeviceUUID string,
) (int, error) {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return 0, ErrInvalidBalenaDeviceUUID
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[DeviceID]{}).
		Get("/v6/device(uuid='" + balenaDeviceUUID + "')?$select=id")
	if err != nil {
		return 0, fmt.Errorf("failed performing request to get device(%s) ID: %w", balenaDeviceUUID, err)
	}

	if response.IsError() {
		return 0, fmt.Errorf("error getting device(%s) ID: %s", balenaDeviceUUID, response.Body())
	}

	balenaResult := response.Result().(*Response[DeviceID])
	if len(balenaResult.D) == 0 {
		return 0, ErrResourceNotFound
	}

	if len(balenaResult.D) > 1 {
		return 0, ErrExpectedOneResult
	}

	return balenaResult.D[0].ID, nil
}

func (b *CloudClient) GetFleet(ctx context.Context, name string) (*Fleet, error) {
	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[Fleet]{}).
		Get("/v6/application?$filter=app_name%20eq%20'" + name + "'")
	if err != nil {
		return nil, fmt.Errorf("failed performing request to get fleet(%s): %w", name, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting fleet(%s): %s", name, response.Body())
	}

	balenaResult := response.Result().(*Response[Fleet])
	if len(balenaResult.D) == 0 {
		return nil, ErrResourceNotFound
	}

	if len(balenaResult.D) > 1 {
		return nil, ErrExpectedOneResult
	}

	return &balenaResult.D[0], nil
}

func (b *CloudClient) RegisterDevice(
	ctx context.Context,
	balenaDeviceUUID, fleetName string,
	deviceType DeviceType,
) error {
	fleet, err := b.GetFleet(ctx, fleetName)
	if err != nil {
		return fmt.Errorf("failed getting fleet(%s): %w", fleetName, err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"application": fleet.ID,
			"uuid":        balenaDeviceUUID,
			"device_type": string(deviceType),
		}).
		Post("/device/register")
	if err != nil {
		return fmt.Errorf("failed performing request to register device(%s) request: %w", balenaDeviceUUID, err)
	}

	if response.IsError() {
		return fmt.Errorf("error while registering device: %s", response.Body())
	}

	return nil
}

func (b *CloudClient) DeleteDevice(
	ctx context.Context,
	balenaDeviceUUID string,
) error {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return ErrInvalidBalenaDeviceUUID
	}

	id, err := b.GetDeviceID(ctx, balenaDeviceUUID)
	if err != nil {
		return fmt.Errorf("failed getting device(%s) ID: %w", balenaDeviceUUID, err)
	}

	response, err := b.httpClient.R().Delete("/v6/device(" + strconv.Itoa(id) + ")")
	if err != nil {
		return fmt.Errorf("failed performing request to delete device(%s): %w", balenaDeviceUUID, err)
	}

	if response.IsError() {
		return fmt.Errorf("error deleting device(%s): %s", balenaDeviceUUID, response.Body())
	}

	return nil
}

func (b *CloudClient) GetDeviceEnvVars(
	ctx context.Context,
	balenaDeviceUUID string,
) ([]DeviceEnvVar, error) {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return nil, ErrInvalidBalenaDeviceUUID
	}

	id, err := b.GetDeviceID(ctx, balenaDeviceUUID)
	if err != nil {
		return nil, fmt.Errorf("failed getting device(%s) ID: %w", balenaDeviceUUID, err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[DeviceEnvVar]{}).
		Get("/v6/device_environment_variable?$filter=device%20eq%20" + strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf("failed performing request for getting device(%s) env vars: %w", balenaDeviceUUID, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting device(%s) env vars: %s", balenaDeviceUUID, response.Body())
	}

	return response.Result().(*Response[DeviceEnvVar]).D, nil
}

func (b *CloudClient) CreateDeviceEnvVar(
	ctx context.Context, balenaDeviceUUID, key string, value interface{},
) error {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return ErrInvalidBalenaDeviceUUID
	}

	id, err := b.GetDeviceID(ctx, balenaDeviceUUID)
	if err != nil {
		return fmt.Errorf("failed getting device(%s) ID: %w", balenaDeviceUUID, err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"device": id,
			"name":   key,
			"value":  value,
		}).
		Post("/v6/device_environment_variable")
	if err != nil {
		return fmt.Errorf("failed performing request to create device(%s) env var(%s): %w", balenaDeviceUUID, key, err)
	}

	if response.IsError() {
		return fmt.Errorf("error creating device(%s) env var(%s): %s", balenaDeviceUUID, key, response.Body())
	}

	return nil
}

func (b *CloudClient) GetFleetEnvVars(
	ctx context.Context,
	name string,
) ([]FleetEnvVar, error) {
	if name == "" {
		return nil, fmt.Errorf("fleet name is required")
	}

	fleet, err := b.GetFleet(ctx, name)
	if err != nil {
		return nil, err
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[FleetEnvVar]{}).
		Get("/v6/application_environment_variable?$filter=application%20eq%20" + strconv.Itoa(fleet.ID))
	if err != nil {
		return nil, fmt.Errorf("failed performing request for getting fleet(%s) env vars: %w", name, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting fleet(%s) env vars: %s", name, response.Body())
	}

	return response.Result().(*Response[FleetEnvVar]).D, nil
}

func (b *CloudClient) GetServiceEnvVars(
	ctx context.Context,
	fleetName string,
) ([]ServiceEnvVar, error) {
	if fleetName == "" {
		return nil, fmt.Errorf("fleet name is required")
	}

	fleet, err := b.GetFleet(ctx, fleetName)
	if err != nil {
		return nil, err
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[ServiceEnvVar]{}).
		Get("/v6/service_environment_variable?$filter=service/any(s:s/application%20eq%20" + strconv.Itoa(fleet.ID) + ")" + "&$select=id,name,value&$expand=service($select=id,service_name)")
	if err != nil {
		return nil, fmt.Errorf("failed performing request for getting service fleet(%s) env vars: %w", fleetName, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting service fleet(%s) env vars: %s", fleetName, response.Body())
	}

	return response.Result().(*Response[ServiceEnvVar]).D, nil
}

func (b *CloudClient) GetDeviceServiceEnvVars(
	ctx context.Context,
	balenaDeviceUUID string,
) ([]DeviceServiceEnvVar, error) {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return nil, ErrInvalidBalenaDeviceUUID
	}

	id, err := b.GetDeviceID(ctx, balenaDeviceUUID)
	if err != nil {
		return nil, fmt.Errorf("failed getting device(%s) ID: %w", balenaDeviceUUID, err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[DeviceServiceEnvVar]{}).
		Get("/v6/device_service_environment_variable?$filter=service_install/any(si:si/device%20eq%20" + strconv.Itoa(id) + ")" + "&$select=id,name,value&$select=id,name,value&$expand=service_install($select=id;$expand=installs__service($select=id,service_name))")
	if err != nil {
		return nil, fmt.Errorf("failed performing request for getting device(%s) service fleet env vars: %w", balenaDeviceUUID, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting device(%s) service fleet env vars: %s", balenaDeviceUUID, response.Body())
	}

	balenaResult := response.Result().(*Response[DeviceServiceEnvVar])
	return balenaResult.D, nil
}

func (b *CloudClient) SetDeviceName(ctx context.Context, balenaDeviceUUID, name string) error {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return ErrInvalidBalenaDeviceUUID
	}

	id, err := b.GetDeviceID(ctx, balenaDeviceUUID)
	if err != nil {
		return fmt.Errorf("failed getting device(%s) ID: %w", balenaDeviceUUID, err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"device_name": name,
		}).
		Patch("/v6/device(" + strconv.Itoa(id) + ")")
	if err != nil {
		return fmt.Errorf("failed performing request to set device(%s) name(%s): %w", balenaDeviceUUID, name, err)
	}

	if response.IsError() {
		return fmt.Errorf("error setting device(%s) name(%s): %s", balenaDeviceUUID, name, response.Body())
	}

	return nil
}

func (b *CloudClient) DownloadOS(
	ctx context.Context, writer io.Writer, fleet string, deviceType DeviceType,
) error {
	flt, err := b.GetFleet(ctx, fleet)
	if err != nil {
		return err
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"deviceType":      string(deviceType),
			"appId":           fmt.Sprintf("%d", flt.ID),
			"fileType":        ".zip",
			"version":         "latest",
			"network":         "ethernet",
			"developmentMode": "false",
		}).
		SetDoNotParseResponse(true).
		Get("/download")
	if err != nil {
		return fmt.Errorf("failed performing request to download os: %w", err)
	}
	defer response.RawResponse.Body.Close()

	if response.IsError() {
		return fmt.Errorf("error downloading os: %s", response.Body())
	}

	_, err = io.Copy(writer, response.RawResponse.Body)
	if err != nil {
		return fmt.Errorf("failed copying response body to writer: %w", err)
	}

	return nil
}

// func (b *CloudClient) ConfigureOSImage(
// 	ctx context.Context,
// 	file, fleet, version string,
// ) error {
// 	cmd := exec.CommandContext(ctx, "balena",
// 		"os", "configure", file,
// 		"--version", version,
// 		"--config-network", "ethernet",
// 		"--fleet", fleet,
// 	)

// 	var dumpOut bytes.Buffer
// 	var dumpErr bytes.Buffer
// 	cmd.Stdout = &dumpOut
// 	cmd.Stderr = &dumpErr

// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println(dumpOut.String())
// 		fmt.Println(dumpErr.String())
// 		return fmt.Errorf("failed configuring os image: %w", err)
// 	}

// 	return nil
// }

func (b *CloudClient) MoveDeviceToFleet(
	ctx context.Context,
	balenaDeviceUUID, fleetName string,
) error {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return ErrInvalidBalenaDeviceUUID
	}

	fleet, err := b.GetFleet(ctx, fleetName)
	if err != nil {
		return err
	}

	deviceID, err := b.GetDeviceID(ctx, balenaDeviceUUID)
	if err != nil {
		return fmt.Errorf("failed getting device(%s) ID: %w", balenaDeviceUUID, err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"belongs_to__application": fleet.ID,
		}).
		Patch("/v6/device(" + strconv.Itoa(deviceID) + ")")
	if err != nil {
		return fmt.Errorf("failed to perform request to move device(%s) to fleet(%s): %w", balenaDeviceUUID, fleetName, err)
	}

	if response.IsError() {
		return fmt.Errorf("error trying to move device(%s) to fleet(%s): %s", balenaDeviceUUID, fleetName, response.Body())
	}

	return nil
}

func (b *CloudClient) EnablePublicDeviceURL(
	ctx context.Context,
	balenaDeviceUUID string,
) error {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return ErrInvalidBalenaDeviceUUID
	}

	id, err := b.GetDeviceID(ctx, balenaDeviceUUID)
	if err != nil {
		return fmt.Errorf("failed getting device(%s) ID: %w", balenaDeviceUUID, err)
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"is_web_accessible": true,
		}).
		Patch("/v6/device(" + strconv.Itoa(id) + ")")
	if err != nil {
		return fmt.Errorf("failed to perform request to enable public device url(%s): %w", balenaDeviceUUID, err)
	}

	if response.IsError() {
		return fmt.Errorf("error trying to enable public device url(%s): %s", balenaDeviceUUID, response.Body())
	}

	return nil
}

func (b *CloudClient) HostLogin(token string) error {
	cmd := exec.Command("balena", "login", "--token", token)

	var dumpOut, dumpErr bytes.Buffer
	cmd.Stdout = &dumpOut
	cmd.Stderr = &dumpErr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to login to Balena: %w", err)
	}

	return nil
}

func (b *CloudClient) GetFleetReleases(
	ctx context.Context,
	name string,
) ([]Release, error) {
	fleet, err := b.GetFleet(ctx, name)
	if err != nil {
		return nil, err
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetResult(Response[Release]{}).
		Get("/v6/release?$filter=belongs_to__application%20eq%20" + strconv.Itoa(fleet.ID) + "&" + OrderByCreatedAtQuerySelector + "&" + FleetReleasesQuerySelect)
	if err != nil {
		return nil, fmt.Errorf("failed performing request to get fleet %s(%d) releases: %w", name, fleet.ID, err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("error getting fleet %s(%d) releases: %s", name, fleet.ID, response.Body())
	}

	balenaResult := response.Result().(*Response[Release])
	if len(balenaResult.D) == 0 {
		return nil, ErrResourceNotFound
	}
	return balenaResult.D, nil
}

func (b *CloudClient) PinDeviceToRelease(
	ctx context.Context,
	balenaDeviceUUID string,
	releaseID int,
) error {
	if !IsValidBalenaDeviceUUID(balenaDeviceUUID) {
		return ErrInvalidBalenaDeviceUUID
	}

	response, err := b.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"should_be_running__release": strconv.Itoa(releaseID),
		}).
		Patch("/v6/device(uuid='" + balenaDeviceUUID + "')")
	if err != nil {
		return fmt.Errorf("failed performing request to pin device(%s) to release(%d): %w", balenaDeviceUUID, releaseID, err)
	}

	if response.IsError() {
		return fmt.Errorf("error trying to pin device(%s) to release(%d): %s", balenaDeviceUUID, releaseID, response.Body())
	}

	return nil
}
