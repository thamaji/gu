package tuple

import (
	"encoding/json"
)

func NewT2[V1, V2 any](v1 V1, v2 V2) T2[V1, V2] {
	return T2[V1, V2]{v1, v2}
}

type T2[V1, V2 any] struct {
	V1 V1
	V2 V2
}

func (t T2[V1, V2]) Values() (V1, V2) {
	return t.V1, t.V2
}

func (t T2[V1, V2]) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(t.V1)
	if err != nil {
		return nil, err
	}

	b2, err := json.Marshal(t.V2)
	if err != nil {
		return nil, err
	}

	return json.Marshal([]json.RawMessage{b1, b2})
}

func (t *T2[V1, V2]) UnmarshalJSON(p []byte) error {
	var msg [2]json.RawMessage
	if err := json.Unmarshal(p, &msg); err != nil {
		return err
	}

	var v1 V1
	if err := json.Unmarshal(msg[0], &v1); err != nil {
		return err
	}

	var v2 V2
	if err := json.Unmarshal(msg[1], &v2); err != nil {
		return err
	}

	t.V1 = v1
	t.V2 = v2

	return nil
}
