// To run this test:
//
// ```bash
// # Note: if you have doppler, you can select the device-manager/local config,
// # and it should have some of these variables set for you.
//
// # Adjust the following environment variables to your own values.
// export BALENA_API_URL=https://api.balena-cloud.com/
// # You can generate an API key from the Balena Cloud dashboard, by going to
// # your personal settings and generating an access token.
// export BALENA_API_KEY=<your-api-key>
// # You will need a device to run this test.
// export BALENA_DEVICE_UUID=<your-device-uuid>
// export TEST_ENV_VAR_NAME=SOME_ENV_VAR_NAME
// # You will need to pick an existing service on the device.
// export EXISTING_SERVICE_NAME=postgres-service
// go run test/e2e/crud_envs.go
//
// ```

//go:build e2e

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Round2POS/gobalena/v2"
)

// Utility function to delete a device env var with a given name, if it exists.
func deleteEnvVarByName(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string) error {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Deleting device env vars: %s", name)
	envVars, err := client.GetDeviceEnvVars(context.Background(), balenaDeviceUUID)
	if err != nil {
		return fmt.Errorf("error getting device env vars: %v", err)
	}

	for _, envVar := range envVars {
		if envVar.Name != name {
			continue
		}

		err := client.DeleteDeviceEnvVar(context.Background(), balenaDeviceID, envVar.ID)
		if err != nil {
			return fmt.Errorf("error deleting device env var: %v", err)
		}
	}

	return nil
}

// Utility function to delete all device service env vars with a given name, if they exist.
func deleteServiceEnvVarsByName(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string) error {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Deleting device service env vars: %s", name)
	serviceEnvVars, err := client.GetDeviceServiceEnvVars(context.Background(), balenaDeviceUUID)
	if err != nil {
		return fmt.Errorf("error getting device service env vars: %v", err)
	}

	for _, serviceEnvVar := range serviceEnvVars {
		if serviceEnvVar.Name != name {
			continue
		}

		err := client.DeleteDeviceServiceEnvVar(context.Background(), balenaDeviceID, serviceEnvVar.ID)
		if err != nil {
			return fmt.Errorf("error deleting device service env var: %v", err)
		}
	}
	return nil
}

// Utility function to get a device env var by name.
func getEnvVarByName(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string) (*gobalena.DeviceEnvVar, error) {
	envVars, err := client.GetDeviceEnvVars(context.Background(), balenaDeviceUUID)
	if err != nil {
		return nil, fmt.Errorf("error getting device env vars: %v", err)
	}

	for _, envVar := range envVars {
		if envVar.Name != name {
			continue
		}

		return &envVar, nil
	}

	return nil, nil
}

// Utility function to get a service install ID by service name.
func getServiceInstallID(
	client gobalena.CloudClient,
	balenaDeviceUUID string,
	serviceName string,
) (int, error) {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Getting service install ID for: %s", serviceName)
	services, err := client.GetDeviceServiceInstallIDs(context.Background(), balenaDeviceUUID)
	if err != nil {
		return 0, fmt.Errorf("error getting service install id: %v", err)
	}

	for _, service := range services {
		if service.ServiceName != serviceName {
			continue
		}

		return service.ServiceInstallID, nil
	}

	return 0, fmt.Errorf("service not found (service name: %s)", serviceName)
}

// Utility function to get a device service env var by name and service install ID.
func getServiceEnvVarByName(client gobalena.CloudClient, balenaDeviceUUID string, name string, serviceInstallID int) (*gobalena.DeviceServiceEnvVar, error) {
	serviceEnvVars, err := client.GetDeviceServiceEnvVars(context.Background(), balenaDeviceUUID)
	if err != nil {
		return nil, fmt.Errorf("error getting device service env vars: %v", err)
	}

	for _, serviceEnvVar := range serviceEnvVars {
		if serviceEnvVar.Name != name {
			continue
		}

		for _, serviceInstall := range serviceEnvVar.ServiceInstall {
			if serviceInstall.ID != serviceInstallID {
				continue
			}

			return &serviceEnvVar, nil
		}
	}

	return nil, nil
}

// This is called at the beginning of each test to clear the device env var and reset the test state.
func clearEnvVar(
	client gobalena.CloudClient,
	balenaDeviceUUID string,
	balenaDeviceID int,
	name string,
) error {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Deleting device env var: %s", name)
	err := deleteEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error deleting device env var: %v", err)
	}
	log.Printf("Device env var deleted: %s", name)
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Confirming deleted device env var: %s", name)
	envVar, err := getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error getting device env var: %v", err)
	}

	if envVar != nil {
		return fmt.Errorf("device env var not deleted: %s", name)
	}

	log.Printf("Confirmed device env var deleted: %s", name)
	//////////////////////////////////////////////////////////////////////////////

	return nil
}

// This is called at the beginning of each test to clear the device service env var and reset the test state.
func clearServiceEnvVar(
	client gobalena.CloudClient,
	balenaDeviceUUID string,
	balenaDeviceID int,
	name string,
	serviceInstallID int,
) error {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Deleting device service env var: %s", name)
	err := deleteServiceEnvVarsByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error deleting device service env var: %v", err)
	}
	log.Printf("Device service env var deleted: %s", name)
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Confirming deleted device service env var: %s", name)
	serviceEnvVar, err := getServiceEnvVarByName(client, balenaDeviceUUID, name, serviceInstallID)
	if err != nil {
		return fmt.Errorf("error getting device service env var: %v", err)
	}

	if serviceEnvVar != nil {
		return fmt.Errorf("device service env var not deleted: %s", name)
	}

	log.Printf("Confirmed device service env var deleted: %s", name)

	return nil
}

// Test deleting a regular device env var.
func testDeleteDeviceEnvVar(
	client gobalena.CloudClient,
	balenaDeviceUUID string,
	balenaDeviceID int,
	name string,
) error {
	err := clearEnvVar(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error clearing device env var: %v", err)
	}

	log.Printf("Creating device env var: %s", name)
	err = client.CreateDeviceEnvVar(context.Background(), balenaDeviceUUID, name, "a value")
	if err != nil {
		return fmt.Errorf("error creating device env var: %v", err)
	}

	log.Printf("Confirming device env var created: %s", name)
	envVar, err := getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error getting device env var: %v", err)
	}

	if envVar == nil {
		return fmt.Errorf("device env var not created: %s", name)
	}

	log.Printf("Deleting device env var: %s", name)
	err = client.DeleteDeviceEnvVar(context.Background(), balenaDeviceID, envVar.ID)
	if err != nil {
		return fmt.Errorf("error deleting device env var: %v", err)
	}

	log.Printf("Confirming device env var deleted: %s", name)
	envVar, err = getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error getting device env var: %v", err)
	}

	if envVar != nil {
		return fmt.Errorf("device env var not deleted: %s", name)
	}

	log.Printf("Confirmed device env var deleted: %s", name)

	return nil
}

func testDeleteDeviceServiceEnvVar(
	client gobalena.CloudClient,
	balenaDeviceUUID string,
	balenaDeviceID int,
	name string,
	serviceName string,
) error {
	serviceInstallID, err := getServiceInstallID(client, balenaDeviceUUID, serviceName)
	if err != nil {
		return fmt.Errorf("error getting service install id: %v", err)
	}

	err = clearServiceEnvVar(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)
	if err != nil {
		return fmt.Errorf("error clearing device service env var: %v", err)
	}

	log.Printf("Creating device service env var: %s", name)
	err = client.CreateDeviceServiceEnvVar(context.Background(), balenaDeviceUUID, name, serviceInstallID, "a value")
	if err != nil {
		return fmt.Errorf("error creating device service env var: %v", err)
	}

	log.Printf("Confirming device service env var created: %s", name)

	serviceEnvVar, err := getServiceEnvVarByName(client, balenaDeviceUUID, name, serviceInstallID)
	if err != nil {
		return fmt.Errorf("error getting device service env var: %v", err)
	}

	if serviceEnvVar == nil {
		return fmt.Errorf("device service env var not created: %s", name)
	}

	log.Printf("Deleting device service env var: %s", name)
	err = client.DeleteDeviceServiceEnvVar(context.Background(), balenaDeviceID, serviceEnvVar.ID)
	if err != nil {
		return fmt.Errorf("error deleting device service env var: %v", err)
	}

	log.Printf("Confirming device service env var deleted: %s", name)
	serviceEnvVar, err = getServiceEnvVarByName(client, balenaDeviceUUID, name, serviceInstallID)
	if err != nil {
		return fmt.Errorf("error getting device service env var: %v", err)
	}

	if serviceEnvVar != nil {
		return fmt.Errorf("device service env var not deleted: %s", name)
	}

	log.Printf("Confirmed device service env var deleted: %s", name)

	return nil
}

func testUpdateDeviceEnvVar(
	client gobalena.CloudClient,
	balenaDeviceUUID string,
	balenaDeviceID int,
	name string,
) error {
	err := clearEnvVar(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error clearing device env var: %v", err)
	}

	log.Printf("Creating device env var: %s", name)
	err = client.CreateDeviceEnvVar(context.Background(), balenaDeviceUUID, name, "a value")
	if err != nil {
		return fmt.Errorf("error creating device env var: %v", err)
	}

	log.Printf("Confirming device env var created: %s", name)
	envVar, err := getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error getting device env var: %v", err)
	}

	log.Printf("Updating device env var: %s", name)
	err = client.UpdateDeviceEnvVar(context.Background(), balenaDeviceID, envVar.ID, "a new value")
	if err != nil {
		return fmt.Errorf("error updating device env var: %v", err)
	}

	log.Printf("Confirming device env var updated: %s", name)
	envVar, err = getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		return fmt.Errorf("error getting device env var: %v", err)
	}

	if envVar.Value != "a new value" {
		return fmt.Errorf("device env var not updated: %s", name)
	}

	log.Printf("Confirmed device env var updated: %s", name)

	return nil
}

func testUpdateDeviceServiceEnvVar(
	client gobalena.CloudClient,
	balenaDeviceUUID string,
	balenaDeviceID int,
	name string,
	serviceName string,
) error {
	serviceInstallID, err := getServiceInstallID(client, balenaDeviceUUID, serviceName)
	if err != nil {
		return fmt.Errorf("error getting service install id: %v", err)
	}

	err = clearServiceEnvVar(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)
	if err != nil {
		return fmt.Errorf("error clearing device service env var: %v", err)
	}

	log.Printf("Creating device service env var: %s", name)
	err = client.CreateDeviceServiceEnvVar(context.Background(), balenaDeviceUUID, name, serviceInstallID, "a value")
	if err != nil {
		return fmt.Errorf("error creating device service env var: %v", err)
	}

	log.Printf("Confirming device service env var created: %s", name)
	serviceEnvVar, err := getServiceEnvVarByName(client, balenaDeviceUUID, name, serviceInstallID)
	if err != nil {
		return fmt.Errorf("error getting device service env var: %v", err)
	}

	log.Printf("Updating device service env var: %s", name)
	err = client.UpdateDeviceServiceEnvVar(context.Background(), balenaDeviceID, serviceEnvVar.ID, "a new value")
	if err != nil {
		return fmt.Errorf("error updating device service env var: %v", err)
	}

	log.Printf("Confirming device service env var updated: %s", name)
	serviceEnvVar, err = getServiceEnvVarByName(client, balenaDeviceUUID, name, serviceInstallID)
	if err != nil {
		return fmt.Errorf("error getting device service env var: %v", err)
	}

	if serviceEnvVar.Value != "a new value" {
		return fmt.Errorf("device service env var not updated: %s", name)
	}

	log.Printf("Confirmed device service env var updated: %s", name)

	return nil
}

func main() {

	apiKey := os.Getenv("BALENA_API_KEY")
	endpoint := os.Getenv("BALENA_API_URL")
	balenaDeviceUUID := os.Getenv("BALENA_DEVICE_UUID")
	existingServiceName := os.Getenv("EXISTING_SERVICE_NAME")
	testEnvVarName := os.Getenv("TEST_ENV_VAR_NAME")

	if apiKey == "" {
		log.Fatalf("Missing environment variables: BALENA_API_KEY")
	}

	if endpoint == "" {
		log.Fatalf("Missing environment variables: BALENA_ENDPOINT")
	}

	if balenaDeviceUUID == "" {
		log.Fatalf("Missing environment variables: BALENA_DEVICE_UUID")
	}

	if testEnvVarName == "" {
		log.Fatalf("Missing environment variables: TEST_ENV_VAR_NAME")
	}

	if existingServiceName == "" {
		log.Fatalf("Missing environment variables: EXISTING_SERVICE_NAME")
	}

	client := gobalena.NewCloudClient(apiKey, endpoint)

	balenaDeviceID, err := client.GetDeviceID(context.Background(), balenaDeviceUUID)
	if err != nil {
		log.Fatalf("Error getting device ID (for device %s): %v", balenaDeviceUUID, err)
	}

	passOrFail := map[string]error{}

	err = testDeleteDeviceEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName)
	passOrFail["testDeleteDeviceEnvVar"] = err

	err = testDeleteDeviceServiceEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName, existingServiceName)
	passOrFail["testDeleteDeviceServiceEnvVar"] = err

	err = testUpdateDeviceEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName)
	passOrFail["testUpdateDeviceEnvVar"] = err

	err = testUpdateDeviceServiceEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName, existingServiceName)
	passOrFail["testUpdateDeviceServiceEnvVar"] = err

	exitCode := 0
	for testName, err := range passOrFail {
		if err != nil {
			log.Printf("❌ Test failed: %s: %v", testName, err)
			exitCode = 1
		} else {
			log.Printf("✅ Test passed: %s", testName)
		}
	}

	err = clearEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName)
	if err != nil {
		log.Fatalf("Error clearing device env var: %v", err)
	}

	log.Printf("Done")
	os.Exit(exitCode)
}
