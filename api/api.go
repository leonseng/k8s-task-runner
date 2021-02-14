package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"k8s.io/client-go/kubernetes"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/leonseng/go_pytest_runner/k8sclient"

	log "github.com/sirupsen/logrus"
)

func HandleRequests(clientset *kubernetes.Clientset, namespace string, port int) {
	r := mux.NewRouter()

	/*
		POST /
		Create a single-run K8s pod (retartPolicy=Never) from the provided image

		Request JSON parameters:
			image: Docker image to run
			command: Overrides command field in the container (equivalent to Docker ENTRYPOINT)
			args: Overrides arguments defined in the container (equivalent to Docker CMD)

		Return 201 if Job was created successfully, 400 otherwise

		Response body:
			id: Request id - used to query for status
			image: Docker image to run
			command: Overrides command field in the container (equivalent to Docker ENTRYPOINT)
			args: Overrides arguments defined in the container (equivalent to Docker CMD)
	*/
	r.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			var reqBody k8sclient.CreateParameters
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			log.Debugf("POST /\n%+v\n-----")
			log.Debugf("", reqBody)

			if err != nil {
				// unable to convert request body to JSON, return 400
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			respBody := reqBody
			respBody.Id = uuid.NewString()

			err = k8sclient.CreatePod(clientset, namespace, respBody)
			if err != nil {
				http.Error(w, "Pod creation has failed:\n"+err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(respBody)
		},
	).Methods(http.MethodPost)

	/*
		GET /{id}
		Gets status of single-run pod, and container logs if test run has been completed

		Path parameter:
			id: Request id

		Return 200 if no errors, 400 otherwise

		Response body:
			id: Request id
			status: Phase of single-run pod
			logs: Terminal output from kubectl logs <pod> command
	*/
	r.HandleFunc(
		"/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			params := mux.Vars(r)
			id := params["id"]
			pod, err := k8sclient.GetPod(clientset, namespace, id)

			if err != nil {
				http.Error(w, "Error getting pod "+id, http.StatusBadRequest)
				return
			}

			respBody := k8sclient.GetStatusParameters{
				Id:     id,
				Status: string(pod.Status.Phase),
			}

			if pod.Status.Phase == "Failed" || pod.Status.Phase == "Succeeded" {
				logs, err := k8sclient.GetPodLogs(clientset, namespace, id)
				if err != nil {
					http.Error(w, "Failed to get pod logs:\n"+err.Error(), http.StatusBadRequest)
					return
				}

				respBody.Logs = logs
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(respBody)
		},
	).Methods(http.MethodGet)

	/*
		GET /{id}/results
		Zip results if not exists, serve zipped results
		Return 200 if no errors
			200 - pod completed
	*/

	http.ListenAndServe(":"+strconv.Itoa(port), r)
}
