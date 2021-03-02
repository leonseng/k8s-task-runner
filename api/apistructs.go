package api

type DockerRegistry struct {
	Server   string `json:"server,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type CreateRequest struct {
	Image          string          `json:"image,omitempty"`
	Command        []string        `json:"command,omitempty"`
	Arguments      []string        `json:"args,omitempty"`
	DockerRegistry *DockerRegistry `json:"dockerRegistry,omitempty"`
}

type CreateResponse struct {
	ID      string        `json:"id"`
	Request CreateRequest `json:"request"`
}

type GetResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Logs   string `json:"logs"`
}
