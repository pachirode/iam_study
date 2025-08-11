package v1

import "encoding/json"

type Response struct {
	Allowed bool   `json:"allowed"`
	Denied  bool   `json:"denied,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (r *Response) String() string {
	data, _ := json.Marshal(r)

	return string(data)
}
