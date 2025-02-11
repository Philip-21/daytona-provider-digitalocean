package provider_test

import (
	"encoding/json"
	"testing"

	"github.com/daytonaio/daytona/pkg/gitprovider"
	daytona_provider "github.com/daytonaio/daytona/pkg/provider"
	"github.com/daytonaio/daytona/pkg/workspace"
	"github.com/daytonaio/daytona/pkg/workspace/project"

	"github.com/daytonaio/daytona-provider-digitalocean/pkg/provider"
	"github.com/daytonaio/daytona-provider-digitalocean/pkg/types"
)

var doProvider = &provider.DigitalOceanProvider{}
var targetOptions = &types.TargetOptions{
	Region:    "nyc3",
	Size:      "s-1vcpu-1gb",
	Image:     "ubuntu-20-04-x64",
	AuthToken: &[]string{"DO_AUTH_TOKEN"}[0],
	DiskSize:  20,
}
var optionsString string

var project1 = &project.Project{
	Name: "test",
	Repository: &gitprovider.GitRepository{
		Id:   "123",
		Url:  "https://github.com/daytonaio/daytona",
		Name: "daytona",
	},
	EnvVars: map[string]string{
		"DAYTONA_WS_ID":                     "123",
		"DAYTONA_WS_PROJECT_NAME":           "test",
		"DAYTONA_WS_PROJECT_REPOSITORY_URL": "https://github.com/daytonaio/daytona",
		"DAYTONA_SERVER_API_KEY":            "api-key-test",
		"DAYTONA_SERVER_VERSION":            "latest",
		"DAYTONA_SERVER_URL":                "http://localhost:3001",
		"DAYTONA_SERVER_API_URL":            "http://localhost:3000",
	},
	WorkspaceId: "123",
}

var workspace1 = &workspace.Workspace{
	Id:     "123",
	Name:   "test",
	Target: "local",
	Projects: []*project.Project{
		project1,
	},
}

func TestCreateWorkspace(t *testing.T) {
	wsReq := &daytona_provider.WorkspaceRequest{
		TargetOptions: optionsString,
		Workspace:     workspace1,
	}

	_, err := doProvider.CreateWorkspace(wsReq)
	if err != nil {
		t.Errorf("Error creating workspace: %s", err)
	}
}

func TestGetWorkspaceInfo(t *testing.T) {
	wsReq := &daytona_provider.WorkspaceRequest{
		TargetOptions: optionsString,
		Workspace:     workspace1,
	}

	workspaceInfo, err := doProvider.GetWorkspaceInfo(wsReq)
	if err != nil || workspaceInfo == nil {
		t.Errorf("Error getting workspace info: %s", err)
	}

	var workspaceMetadata types.WorkspaceMetadata
	err = json.Unmarshal([]byte(workspaceInfo.ProviderMetadata), &workspaceMetadata)
	if err != nil {
		t.Errorf("Error unmarshalling workspace metadata: %s", err)
	}
}

func TestDestroyWorkspace(t *testing.T) {
	wsReq := &daytona_provider.WorkspaceRequest{
		TargetOptions: optionsString,
		Workspace:     workspace1,
	}

	_, err := doProvider.DestroyWorkspace(wsReq)
	if err != nil {
		t.Errorf("Error deleting workspace: %s", err)
	}
}

func TestCreateProject(t *testing.T) {
	TestCreateWorkspace(t)

	projectReq := &daytona_provider.ProjectRequest{
		TargetOptions: optionsString,
		Project:       project1,
	}

	_, err := doProvider.CreateProject(projectReq)
	if err != nil {
		t.Errorf("Error creating project: %s", err)
	}
}

func TestStartProject(t *testing.T) {
	projectReq := &daytona_provider.ProjectRequest{
		TargetOptions: optionsString,
		Project:       project1,
	}

	// Call StartProject
	_, err := doProvider.StartProject(projectReq)
	if err != nil {
		t.Errorf("Error starting a project: %s", err)
	}
}

func TestStopProject(t *testing.T) {
	projectReq := &daytona_provider.ProjectRequest{
		TargetOptions: optionsString,
		Project:       project1,
	}

	// Call StartProject
	_, err := doProvider.StopProject(projectReq)
	if err != nil {
		t.Errorf("Error stopping a project: %s", err)
	}
}

func TestDestroyProject(t *testing.T) {
	projectReq := &daytona_provider.ProjectRequest{
		TargetOptions: optionsString,
		Project:       project1,
	}

	_, err := doProvider.DestroyProject(projectReq)
	if err != nil {
		t.Errorf("Error deleting project: %s", err)
	}

	TestDestroyWorkspace(t)
}

func init() {
	_, err := doProvider.Initialize(daytona_provider.InitializeProviderRequest{
		BasePath:           "/tmp/workspaces",
		DaytonaDownloadUrl: "https://download.daytona.io/daytona/get-server.sh",
		DaytonaVersion:     "latest",
		ServerUrl:          "",
		ApiUrl:             "",
		ServerPort:         0,
		ApiPort:            0,
		LogsDir:            "/tmp/logs",
	})
	if err != nil {
		panic(err)
	}

	opts, err := json.Marshal(targetOptions)
	if err != nil {
		panic(err)
	}

	optionsString = string(opts)
}
