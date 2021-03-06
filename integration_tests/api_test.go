package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/leonseng/k8s_task_runner/api"
	"github.com/stretchr/testify/assert"
)

func TestInCluster(t *testing.T) {
	runIntegrationTest(t, "http://localhost:8080/")
}

func TestOutOfCluster(t *testing.T) {
	runIntegrationTest(t, "http://localhost:8081/")
}

func runIntegrationTest(t *testing.T, apiEndpoint string) {
	getStatus(t, apiEndpoint)
	createAndGetTask(t, apiEndpoint)
}

func getStatus(t *testing.T, apiEndpoint string) {
	resp, err := http.Get(apiEndpoint + "/status")
	if err != nil {
		t.Errorf("Failed to get app status: %v\n", err)
	}

	defer resp.Body.Close()
	respBody := new(api.GetStatusResponse)
	fmt.Printf("%+v\n", resp.Body)
	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		t.Errorf("Failed to decode GET response body to JSON")
	}

	assert.Equal(t, respBody.Status, "healthy")
}

func createAndGetTask(t *testing.T, apiEndpoint string) {
	test_env_var := "test_var"
	test_env_var_value := "123"
	reqBody := api.CreateRequest{
		Image:     "busybox:1.28",
		Command:   []string{"printenv"},
		Arguments: []string{test_env_var},
		DockerRegistry: &api.DockerRegistry{
			Server:   "test.com",
			Username: "test-user",
			Password: "secure",
			Email:    "test-user@test.com",
		},
		EnvVars: map[string]string{
			test_env_var: test_env_var_value,
		},
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(reqBody)
	if err != nil {
		t.Errorf("Failed to encode POST request body to JSON")
	}

	createResp, err := http.Post(apiEndpoint, "application/json", b)
	if err != nil {
		t.Errorf("Failed to create job: %v\n", err)
	}

	assert.Equal(t, 201, createResp.StatusCode)

	defer createResp.Body.Close()
	createRespBody := new(api.CreateResponse)
	fmt.Printf("%+v\n", createResp.Body)
	err = json.NewDecoder(createResp.Body).Decode(createRespBody)
	fmt.Printf("%+v\n", createRespBody)
	if err != nil {
		t.Errorf("Failed to decode POST response body to struct")
	}

	// wait for pod to run to completion
	getRespBody := new(api.GetTaskResponse)
	var getResp *http.Response
	for i := 0; i < 30; i++ {
		getResp, err = http.Get(apiEndpoint + createRespBody.ID)
		if err != nil {
			t.Errorf("Failed to get job status: %v\n", err)
		}

		defer getResp.Body.Close()
		err = json.NewDecoder(getResp.Body).Decode(getRespBody)
		if err != nil {
			t.Errorf("Failed to decode GET response body to struct")
		} else if getRespBody.Status == "Succeeded" {
			break
		}

		time.Sleep(time.Second)
	}

	if getRespBody.Status != "Succeeded" {
		t.Errorf("Test pod failed to run to completion.\nResponse Status Code: %d\nBody:\n%+v", getResp.StatusCode, *getRespBody)
	}

	if strings.TrimSpace(getRespBody.Logs) != test_env_var_value {
		t.Errorf("Test pod failed to read environment variable correctly. Test pod log: \n%+v", getRespBody.Logs)
	}
}
