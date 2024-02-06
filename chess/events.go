package chess

import (
	"encoding/json"
)

type inputEvent struct {
	Client *Client `json:"-"`
	Square string  `json:"HX-Trigger"`
	Action string  `json:"HX-Trigger-Name"`
	URL    string  `json:"HX-Current-URL"`
}

func UnmarshalAction(bytes []byte) (*inputEvent, error) {
	var tmp struct {
		Headers inputEvent `json:"HEADERS"`
	}

	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return nil, err
	}

	return &tmp.Headers, nil
}
