package response

type Err struct {
	Code    int         `json:"-"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload,omitempty"`
}

func (e *Err) Error() string {
	return e.Message
}

type Res struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload,omitempty"`
}
