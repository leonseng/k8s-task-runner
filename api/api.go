package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/leonseng/k8s_task_runner/k8sclient"

	log "github.com/sirupsen/logrus"
)

func HandleRequests(appConf ApplicationConfiguration) {
	r := mux.NewRouter()
	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/swaggerui/").Handler(sh)
	r.HandleFunc("/status", getStatus).Methods(http.MethodGet)
	r.HandleFunc("/task/{id}", getTask(appConf)).Methods(http.MethodGet)
	r.HandleFunc("/task", createTask(appConf)).Methods(http.MethodPost)

	err := http.ListenAndServe(":"+strconv.Itoa(appConf.Port), r)
	if err != nil {
		panic(err.Error())
	}
}

/*
	POST /
	Create a single-run K8s pod (retartPolicy=Never) from the provided image

	Request JSON parameters:
		image: Docker image to run
		command: Overrides command field in the container (equivalent to Docker ENTRYPOINT)
		args: Overrides arguments defined in the container (equivalent to Docker CMD)
		dockerRegistry:
			server: Private Docker Registry FQDN. Use https://index.docker.io/v2/ for DockerHub.
			username: Docker username
			password: Docker password
			email: Docker email

	Return 201 if Job was created successfully, 400 otherwise

	Response body:
		id: Request id - used to query for status
		request: Request JSON body as parsed by k8s-task-runner
*/
func createTask(appConf ApplicationConfiguration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody CreateRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		log.Debugf("POST /\n%+v\n-----", reqBody)

		if err != nil {
			// unable to convert request body to JSON, return 400
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		id := uuid.NewString()

		// create secret if provided
		var secretName string
		if reqBody.DockerRegistry != nil {
			secretName, err = k8sclient.CreateDockerRegistrySecret(
				appConf.K8sClientSet,
				k8sclient.SecretParameters{
					ID:        id,
					Namespace: appConf.TaskNamespace,
					Server:    reqBody.DockerRegistry.Server,
					Username:  reqBody.DockerRegistry.Username,
					Email:     reqBody.DockerRegistry.Email,
					Password:  reqBody.DockerRegistry.Password,
				},
			)
			if err != nil {
				http.Error(w, "Docker registry secret creation has failed:\n"+err.Error(), http.StatusBadRequest)
				return
			}
		}

		err = k8sclient.CreatePodFromManifest(
			appConf.K8sClientSet,
			k8sclient.PodParameters{
				ID:        id,
				Namespace: appConf.TaskNamespace,
				Secret:    secretName,
				Image:     reqBody.Image,
				Command:   reqBody.Command,
				Arguments: reqBody.Arguments,
				EnvVars:   reqBody.EnvVars,
			},
		)

		if err != nil {
			http.Error(w, "Pod creation has failed:\n"+err.Error(), http.StatusBadRequest)
			return
		}

		respData, err := json.Marshal(
			CreateResponse{
				ID:      id,
				Request: reqBody,
			},
		)
		if err != nil {
			log.Error("Failed to encode POST response struct to JSON")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(respData)
	}
}

/*
	GET /status
	Gets status of the k8s-task-runner application

	Return 200

	Response body:
		status: Returns health of application: healthy
*/
func getStatus(w http.ResponseWriter, r *http.Request) {
	respData, err := json.Marshal(GetStatusResponse{Status: "healthy"})
	if err != nil {
		log.Error("Failed to encode GET app status response struct to JSON")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}

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
func getTask(appConf ApplicationConfiguration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		pod, err := k8sclient.GetPod(appConf.K8sClientSet, appConf.TaskNamespace, id)

		if err != nil {
			http.Error(w, "Error getting pod "+id, http.StatusBadRequest)
			return
		}

		respBody := GetTaskResponse{
			ID:     id,
			Status: string(pod.Status.Phase),
		}

		if pod.Status.Phase == "Failed" || pod.Status.Phase == "Succeeded" {
			logs, err := k8sclient.GetPodLogs(appConf.K8sClientSet, appConf.TaskNamespace, id)
			if err != nil {
				http.Error(w, "Failed to get pod logs:\n"+err.Error(), http.StatusBadRequest)
				return
			}

			respBody.Logs = logs
		}

		respData, err := json.Marshal(respBody)
		if err != nil {
			log.Error("Failed to encode GET task status response struct to JSON")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respData)
	}
}
