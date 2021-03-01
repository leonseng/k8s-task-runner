package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/leonseng/k8s_task_runner/api"
	"github.com/stretchr/testify/assert"
)

func TestK8sTestRunner(t *testing.T) {
	reqBody := api.CreateRequest{
		Image:   "busybox:1.28",
		Command: []string{"date"},
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(reqBody)
	if err != nil {
		t.Errorf("Failed to encode POST request body to JSON")
	}

	createResp, err := http.Post("http://localhost:8080", "application/json", b)
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
	getRespBody := new(api.GetResponse)
	var getResp *http.Response
	for i := 0; i < 30; i++ {
		getResp, err = http.Get("http://localhost:8080/" + createRespBody.ID)
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
}
