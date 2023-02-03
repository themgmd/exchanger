package response

type DefaultHttpResponse struct {
	Success bool   `json:"success"`
	Comment string `json:"comment"`
}

type PingPongResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}
