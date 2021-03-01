package api

type CreateRequest struct {
	Image     string   `json:"image"`
	Command   []string `json:"command"`
	Arguments []string `json:"args"`
}

type CreateResponse struct {
	ID  string `json:"id"`
	Request CreateRequest `json:"request"`
}

type GetResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Logs   string `json:"logs"`
}
