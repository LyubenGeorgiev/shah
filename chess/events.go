package chess

import (
	"encoding/json"
)

type inputEvent struct {
	Client  *Client `json:"-"`
	Square  string  `json:"HX-Trigger"`
	Action  string  `json:"HX-Trigger-Name"`
	URL     string  `json:"HX-Current-URL"`
	Message string  `json:"-"`
}

func UnmarshalAction(bytes []byte) (*inputEvent, error) {
	var tmp struct {
		Message string     `json:"chat_message"`
		Headers inputEvent `json:"HEADERS"`
	}

	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return nil, err
	}

	tmp.Headers.Message = tmp.Message

	return &tmp.Headers, nil
}
