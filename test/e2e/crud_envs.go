// To run thsi test:
//
// ```bash
// doppler setup
// # Choose "local" config, "device-manager" project.
//
// export BALENA_DEVICE_UUID=<your-device-uuid>
// export TEST_ENV_VAR_NAME=SOME_ENV_VAR_NAME
// export EXISTING_SERVICE_NAME=postgres-service
// doppler run  -- go run test/e2e/crud_envs.go
// ```

//go:build e2e

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Round2POS/gobalena"
)

func deleteEnvVarsByName(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string) error {
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

func getServiceInstallIDOrDie(client gobalena.CloudClient, balenaDeviceUUID string, serviceName string) int {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Getting service install ID for: %s", serviceName)
	services, err := client.GetDeviceServiceInstallIDs(context.Background(), balenaDeviceUUID)
	if err != nil {
		log.Fatalf("error getting service install id: error getting device service env vars: %v", err)
	}

	for _, service := range services {
		if service.ServiceName != serviceName {
			continue
		}

		return service.ServiceInstallID
	}

	log.Fatalf("error getting service install id: service not found (service name: %s)", serviceName)

	return 0
}

func getServiceIDOrDie(client gobalena.CloudClient, balenaDeviceUUID string, serviceName string) int {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Getting service ID for: %s", serviceName)
	services, err := client.GetDeviceServiceInstallIDs(context.Background(), balenaDeviceUUID)
	if err != nil {
		log.Fatalf("error getting service id: error getting device service env vars: %v", err)
	}

	for _, service := range services {
		if service.ServiceName != serviceName {
			continue
		}

		return service.ServiceID
	}

	log.Fatalf("error getting service id: service not found (service name: %s)", serviceName)

	return 0
}

func getServiceEnvVarByName(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string, serviceInstallID int) (*gobalena.DeviceServiceEnvVar, error) {
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

func clearEnvVarOrDie(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string) {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Deleting device env var: %s", name)
	err := deleteEnvVarsByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		log.Fatalf("Error deleting device env var: %v", err)
	}
	log.Printf("Device env var deleted: %s", name)
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Confirming deleted device env var: %s", name)
	envVar, err := getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		log.Fatalf("Error getting device env var: %v", err)
	}

	if envVar != nil {
		log.Fatalf("Device env var not deleted: %s", name)
	}

	log.Printf("Confirmed device env var deleted: %s", name)
	//////////////////////////////////////////////////////////////////////////////

}

func clearServiceEnvVarOrDie(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string, serviceInstallID int) {
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Deleting device service env var: %s", name)
	err := deleteServiceEnvVarsByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		log.Fatalf("Error deleting device service env var: %v", err)
	}
	log.Printf("Device service env var deleted: %s", name)
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Confirming deleted device service env var: %s", name)
	serviceEnvVar, err := getServiceEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)
	if err != nil {
		log.Fatalf("Error getting device service env var: %v", err)
	}

	if serviceEnvVar != nil {
		log.Fatalf("Device service env var not deleted: %s", name)
	}

	log.Printf("Confirmed device service env var deleted: %s", name)
}

func testDeleteDeviceEnvVar(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string) {
	clearEnvVarOrDie(client, balenaDeviceUUID, balenaDeviceID, name)

	log.Printf("Creating device env var: %s", name)
	err := client.CreateDeviceEnvVar(context.Background(), balenaDeviceUUID, name, "a value")
	if err != nil {
		log.Fatalf("Error creating device env var: %v", err)
	}

	log.Printf("Confirming device env var created: %s", name)
	envVar, err := getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		log.Fatalf("Error getting device env var: %v", err)
	}

	if envVar == nil {
		log.Fatalf("Device env var not created: %s", name)
	}

	log.Printf("Deleting device env var: %s", name)
	err = client.DeleteDeviceEnvVar(context.Background(), balenaDeviceID, envVar.ID)
	if err != nil {
		log.Fatalf("Error deleting device env var: %v", err)
	}

	log.Printf("Confirming device env var deleted: %s", name)
	envVar, err = getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		log.Fatalf("Error getting device env var: %v", err)
	}

	if envVar != nil {
		log.Fatalf("Device env var not deleted: %s", name)
	}

	log.Printf("Confirmed device env var deleted: %s", name)
}

func testDeleteDeviceServiceEnvVar(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string, serviceName string) {

	serviceInstallID := getServiceInstallIDOrDie(client, balenaDeviceUUID, serviceName)
	clearServiceEnvVarOrDie(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)

	log.Printf("Creating device service env var: %s", name)
	err := client.CreateDeviceServiceEnvVar(context.Background(), balenaDeviceUUID, name, serviceInstallID, "a value")
	if err != nil {
		log.Fatalf("Error creating device service env var: %v", err)
	}

	log.Printf("Confirming device service env var created: %s", name)

	serviceEnvVar, err := getServiceEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)
	if err != nil {
		log.Fatalf("Error getting device service env var: %v", err)
	}

	if serviceEnvVar == nil {
		log.Fatalf("Device service env var not created: %s", name)
	}

	log.Printf("Deleting device service env var: %s", name)
	err = client.DeleteDeviceServiceEnvVar(context.Background(), balenaDeviceID, serviceEnvVar.ID)
	if err != nil {
		log.Fatalf("Error deleting device service env var: %v", err)
	}

	log.Printf("Confirming device service env var deleted: %s", name)
	serviceEnvVar, err = getServiceEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)
	if err != nil {
		log.Fatalf("Error getting device service env var: %v", err)
	}

	if serviceEnvVar != nil {
		log.Fatalf("Device service env var not deleted: %s", name)
	}

	log.Printf("Confirmed device service env var deleted: %s", name)
}

func testUpdateDeviceEnvVar(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string) {
	clearEnvVarOrDie(client, balenaDeviceUUID, balenaDeviceID, name)

	log.Printf("Creating device env var: %s", name)
	err := client.CreateDeviceEnvVar(context.Background(), balenaDeviceUUID, name, "a value")
	if err != nil {
		log.Fatalf("Error creating device env var: %v", err)
	}

	log.Printf("Confirming device env var created: %s", name)
	envVar, err := getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		log.Fatalf("Error getting device env var: %v", err)
	}

	log.Printf("Updating device env var: %s", name)
	err = client.UpdateDeviceEnvVar(context.Background(), balenaDeviceID, envVar.ID, "a new value")
	if err != nil {
		log.Fatalf("Error updating device env var: %v", err)
	}

	log.Printf("Confirming device env var updated: %s", name)
	envVar, err = getEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name)
	if err != nil {
		log.Fatalf("Error getting device env var: %v", err)
	}

	if envVar.Value != "a new value" {
		log.Fatalf("Device env var not updated: %s", name)
	}

	log.Printf("Confirmed device env var updated: %s", name)
}

func testUpdateDeviceServiceEnvVar(client gobalena.CloudClient, balenaDeviceUUID string, balenaDeviceID int, name string, serviceName string) {
	serviceInstallID := getServiceInstallIDOrDie(client, balenaDeviceUUID, serviceName)
	clearServiceEnvVarOrDie(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)

	log.Printf("Creating device service env var: %s", name)
	err := client.CreateDeviceServiceEnvVar(context.Background(), balenaDeviceUUID, name, serviceInstallID, "a value")
	if err != nil {
		log.Fatalf("Error creating device service env var: %v", err)
	}

	log.Printf("Confirming device service env var created: %s", name)
	serviceEnvVar, err := getServiceEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)
	if err != nil {
		log.Fatalf("Error getting device service env var: %v", err)
	}

	log.Printf("Updating device service env var: %s", name)
	err = client.UpdateDeviceServiceEnvVar(context.Background(), balenaDeviceID, serviceEnvVar.ID, "a new value")
	if err != nil {
		log.Fatalf("Error updating device service env var: %v", err)
	}

	log.Printf("Confirming device service env var updated: %s", name)
	serviceEnvVar, err = getServiceEnvVarByName(client, balenaDeviceUUID, balenaDeviceID, name, serviceInstallID)
	if err != nil {
		log.Fatalf("Error getting device service env var: %v", err)
	}

	if serviceEnvVar.Value != "a new value" {
		log.Fatalf("Device service env var not updated: %s", name)
	}

	log.Printf("Confirmed device service env var updated: %s", name)
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

	clearEnvVarOrDie(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName)

	testDeleteDeviceEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName)
	testDeleteDeviceServiceEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName, existingServiceName)
	testUpdateDeviceEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName)
	testUpdateDeviceServiceEnvVar(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName, existingServiceName)

	clearEnvVarOrDie(client, balenaDeviceUUID, balenaDeviceID, testEnvVarName)

	log.Printf("Done")
}
