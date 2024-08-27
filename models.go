package gobalena

import "time"

type DeviceType string

type serializableResponse interface {
	Device | DeviceTag | Fleet | EnvVar
}

type Response[T serializableResponse] struct {
	D []T `json:"d"`
}

type Device struct {
	ID                    int       `json:"id"`
	IPAddress             string    `json:"ip_address"`
	MacAddress            string    `json:"mac_address"`
	PublicAddress         string    `json:"public_address"`
	DeviceName            string    `json:"device_name"`
	OsVersion             string    `json:"os_version"`
	OsVariant             string    `json:"os_variant"`
	SupervisorVersion     string    `json:"supervisor_version"`
	IsOnline              bool      `json:"is_online"`
	LastConnectivityEvent time.Time `json:"last_connectivity_event"`
	IsWebAccessible       bool      `json:"is_web_accessible"`
	Latitude              string    `json:"latitude"`
	Longitude             string    `json:"longitude"`
	Location              string    `json:"location"`
	CreatedAt             time.Time `json:"created_at"`
	IsRunningRelease      []struct {
		ID         int    `json:"id"`
		RawVersion string `json:"raw_version"`
		Commit     string `json:"commit"`
	} `json:"is_running__release"`
	ShouldBeRunningRelease []struct {
		ID         int    `json:"id"`
		RawVersion string `json:"raw_version"`
		Commit     string `json:"commit"`
	} `json:"should_be_running__release"`
}

type DeviceTag struct {
	ID     int `json:"id"`
	Device struct {
		ID int `json:"__id"`
	} `json:"device"`
	TagKey string `json:"tag_key"`
	Value  string `json:"value"`
}

type Fleet struct {
	ID           int `json:"id"`
	Organization struct {
		ID int `json:"__id"`
	} `json:"organization"`
	Actor                  int    `json:"actor"`
	AppName                string `json:"app_name"`
	Slug                   string `json:"slug"`
	ShouldBeRunningRelease any    `json:"should_be_running__release"`
	ApplicationType        struct {
		ID int `json:"__id"`
	} `json:"application_type"`
	IsForDeviceType struct {
		ID int `json:"__id"`
	} `json:"is_for__device_type"`
	ShouldTrackLatestRelease       bool      `json:"should_track_latest_release"`
	IsAccessibleBySupportUntilDate any       `json:"is_accessible_by_support_until__date"`
	IsPublic                       bool      `json:"is_public"`
	IsHost                         bool      `json:"is_host"`
	IsArchived                     bool      `json:"is_archived"`
	IsDiscoverable                 bool      `json:"is_discoverable"`
	IsStoredAtRepositoryURL        any       `json:"is_stored_at__repository_url"`
	CreatedAt                      time.Time `json:"created_at"`
	UUID                           string    `json:"uuid"`
	IsOfClass                      string    `json:"is_of__class"`
}

type EnvVar struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Device    struct {
		ID int `json:"__id"`
	} `json:"device"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeviceTags struct {
	Status string `json:"status"`
	Tags   []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"tags"`
}

type Status struct {
	Status                  string `json:"status"`
	AppState                string `json:"appState"`
	OverallDownloadProgress any    `json:"overallDownloadProgress"`
	Containers              []struct {
		Status      string    `json:"status"`
		ServiceName string    `json:"serviceName"`
		AppID       int       `json:"appId"`
		ImageID     int       `json:"imageId"`
		ServiceID   int       `json:"serviceId"`
		ContainerID string    `json:"containerId"`
		CreatedAt   time.Time `json:"createdAt"`
	} `json:"containers"`
	Images []struct {
		Name             string `json:"name"`
		AppID            int    `json:"appId"`
		ServiceName      string `json:"serviceName"`
		ImageID          int    `json:"imageId"`
		DockerImageID    string `json:"dockerImageId"`
		Status           string `json:"status"`
		DownloadProgress any    `json:"downloadProgress"`
	} `json:"images"`
	Release string `json:"release"`
}

type DeviceState struct {
	APIPort           int    `json:"api_port"`
	IPAddress         string `json:"ip_address"`
	OsVersion         string `json:"os_version"`
	MacAddress        string `json:"mac_address"`
	SupervisorVersion string `json:"supervisor_version"`
	UpdatePending     bool   `json:"update_pending"`
	UpdateFailed      bool   `json:"update_failed"`
	UpdateDownloaded  bool   `json:"update_downloaded"`
	Commit            string `json:"commit"`
	Status            string `json:"status"`
	DownloadProgress  any    `json:"download_progress"`
}

type balenaSerializableResponse interface {
	BalenaDevice | BalenaDeviceTag | BalenaFleet |
		BalenaDeviceEnvVar | BalenaRelease | BalenaFleetEnvVar | BalenaServiceEnvVar | BalenaDeviceID | BalenaDeviceServiceEnvVar
}

type BalenaResponse[T balenaSerializableResponse] struct {
	D []T `json:"d"`
}

type BalenaDevice struct {
	ID                     int                `json:"id"`
	UUID                   string             `json:"uuid"`
	IPAddress              string             `json:"ip_address"`
	MacAddress             string             `json:"mac_address"`
	PublicAddress          string             `json:"public_address"`
	DeviceName             string             `json:"device_name"`
	OsVersion              string             `json:"os_version"`
	OsVariant              string             `json:"os_variant"`
	SupervisorVersion      string             `json:"supervisor_version"`
	IsOnline               bool               `json:"is_online"`
	LastConnectivityEvent  time.Time          `json:"last_connectivity_event"`
	IsWebAccessible        bool               `json:"is_web_accessible"`
	Latitude               string             `json:"latitude"`
	Longitude              string             `json:"longitude"`
	Location               string             `json:"location"`
	CreatedAt              time.Time          `json:"created_at"`
	IsRunningRelease       []BalenaRelease    `json:"is_running__release"`
	ShouldBeRunningRelease []BalenaRelease    `json:"should_be_running__release"`
	BelongsToApplication   []BalenaFleetShort `json:"belongs_to__application"`
	OverallStatus          string             `json:"overall_status"`
}

type BalenaDeviceID struct {
	ID int `json:"id"`
}

type BalenaDeviceTag struct {
	ID     int `json:"id"`
	Device struct {
		ID int `json:"__id"`
	} `json:"device"`
	TagKey string `json:"tag_key"`
	Value  string `json:"value"`
}

type BalenaFleet struct {
	ID           int `json:"id"`
	Organization struct {
		ID int `json:"__id"`
	} `json:"organization"`
	Actor                  int    `json:"actor"`
	AppName                string `json:"app_name"`
	Slug                   string `json:"slug"`
	ShouldBeRunningRelease any    `json:"should_be_running__release"`
	ApplicationType        struct {
		ID int `json:"__id"`
	} `json:"application_type"`
	IsForDeviceType struct {
		ID int `json:"__id"`
	} `json:"is_for__device_type"`
	ShouldTrackLatestRelease       bool      `json:"should_track_latest_release"`
	IsAccessibleBySupportUntilDate any       `json:"is_accessible_by_support_until__date"`
	IsPublic                       bool      `json:"is_public"`
	IsHost                         bool      `json:"is_host"`
	IsArchived                     bool      `json:"is_archived"`
	IsDiscoverable                 bool      `json:"is_discoverable"`
	IsStoredAtRepositoryURL        any       `json:"is_stored_at__repository_url"`
	CreatedAt                      time.Time `json:"created_at"`
	UUID                           string    `json:"uuid"`
	IsOfClass                      string    `json:"is_of__class"`
}

type BalenaFleetShort struct {
	ID      int    `json:"id"`
	AppName string `json:"app_name"`
}

type BalenaRelease struct {
	IsCreatedByUser []struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"is_created_by__user"`
	IsRunningOnDevice int `json:"is_running_on__device"`
	ReleaseTag        []struct {
		TagKey string `json:"tag_key"`
		Value  string `json:"value"`
		ID     int    `json:"id"`
	} `json:"release_tag"`
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// type BalenaReleaseLong struct {
// 	IsCreatedByUser []struct {
// 		ID        int       `json:"id"`
// 		Username  string    `json:"username"`
// 		CreatedAt time.Time `json:"created_at"`
// 	} `json:"is_created_by__user"`
// 	IsRunningOnDevice int `json:"is_running_on__device"`
// 	ReleaseTag        []struct {
// 		TagKey string `json:"tag_key"`
// 		Value  string `json:"value"`
// 		ID     int    `json:"id"`
// 	} `json:"release_tag"`
// 	ID                   int `json:"id"`
// 	BelongsToApplication struct {
// 		ID int `json:"__id"`
// 	} `json:"belongs_to__application"`
// 	BuildLog    any    `json:"build_log"`
// 	Commit      string `json:"commit"`
// 	Composition struct {
// 		Version  string      `json:"version"`
// 		Volumes  interface{} `json:"volumes"`
// 		Services interface{} `json:"services"`
// 	} `json:"composition"`
// 	Contract           any       `json:"contract"`
// 	CreatedAt          time.Time `json:"created_at"`
// 	EndTimestamp       time.Time `json:"end_timestamp"`
// 	InvalidationReason any       `json:"invalidation_reason"`
// 	IsFinal            bool      `json:"is_final"`
// 	IsFinalizedAtDate  time.Time `json:"is_finalized_at__date"`
// 	IsInvalidated      bool      `json:"is_invalidated"`
// 	IsPassingTests     bool      `json:"is_passing_tests"`
// 	KnownIssueList     any       `json:"known_issue_list"`
// 	Note               any       `json:"note"`
// 	Phase              any       `json:"phase"`
// 	ReleaseVersion     any       `json:"release_version"`
// 	Revision           int       `json:"revision"`
// 	Semver             string    `json:"semver"`
// 	SemverBuild        string    `json:"semver_build"`
// 	SemverMajor        int       `json:"semver_major"`
// 	SemverMinor        int       `json:"semver_minor"`
// 	SemverPatch        int       `json:"semver_patch"`
// 	SemverPrerelease   string    `json:"semver_prerelease"`
// 	Source             string    `json:"source"`
// 	StartTimestamp     time.Time `json:"start_timestamp"`
// 	Status             string    `json:"status"`
// 	UpdateTimestamp    time.Time `json:"update_timestamp"`
// 	Variant            string    `json:"variant"`
// 	ReleaseType        string    `json:"release_type"`
// }

type BalenaDeviceEnvVar struct {
	ID     int `json:"id"`
	Device struct {
		ID int `json:"__id"`
	} `json:"device"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BalenaFleetEnvVar struct {
	ID          int `json:"id"`
	Application struct {
		ID int `json:"__id"`
	} `json:"application"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BalenaServiceShort struct {
	ID          int    `json:"id"`
	ServiceName string `json:"service_name"`
}

type BalenaServiceEnvVar struct {
	Service []BalenaServiceShort `json:"service"`
	ID      int                  `json:"id"`
	Name    string               `json:"name"`
	Value   string               `json:"value"`
}

// type BalenaDeviceServiceEnvVar struct {
// 	ServiceInstall []struct {
// 		InstallsService []BalenaServiceEnvVar `json:"installs__service"`
// 		ID              int                   `json:"id"`
// 	} `json:"service_install"`
// 	ID    int    `json:"id"`
// 	Name  string `json:"name"`
// 	Value string `json:"value"`
// }

type BalenaDeviceServiceEnvVar struct {
	ServiceInstall []struct {
		InstallsService []BalenaServiceShort `json:"installs__service"`
		ID              int                  `json:"id"`
	} `json:"service_install"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BalenaGenericEnvVar struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	FleetValue  string             `json:"fleet_value"`
	DeviceValue string             `json:"device_value"`
	Service     BalenaServiceShort `json:"service"`
}
