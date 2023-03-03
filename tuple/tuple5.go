package tuple

import "encoding/json"

func NewT5[V1, V2, V3, V4, V5 any](v1 V1, v2 V2, v3 V3, v4 V4, v5 V5) T5[V1, V2, V3, V4, V5] {
	return T5[V1, V2, V3, V4, V5]{v1, v2, v3, v4, v5}
}

type T5[V1, V2, V3, V4, V5 any] struct {
	V1 V1
	V2 V2
	V3 V3
	V4 V4
	V5 V5
}

func (t T5[V1, V2, V3, V4, V5]) Values() (V1, V2, V3, V4, V5) {
	return t.V1, t.V2, t.V3, t.V4, t.V5
}

func (t T5[V1, V2, V3, V4, V5]) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(t.V1)
	if err != nil {
		return nil, err
	}

	b2, err := json.Marshal(t.V2)
	if err != nil {
		return nil, err
	}

	b3, err := json.Marshal(t.V3)
	if err != nil {
		return nil, err
	}

	b4, err := json.Marshal(t.V4)
	if err != nil {
		return nil, err
	}

	b5, err := json.Marshal(t.V5)
	if err != nil {
		return nil, err
	}

	return json.Marshal([]json.RawMessage{b1, b2, b3, b4, b5})
}

func (t *T5[V1, V2, V3, V4, V5]) UnmarshalJSON(p []byte) error {
	var msg [5]json.RawMessage
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

	var v3 V3
	if err := json.Unmarshal(msg[2], &v3); err != nil {
		return err
	}

	var v4 V4
	if err := json.Unmarshal(msg[3], &v4); err != nil {
		return err
	}

	var v5 V5
	if err := json.Unmarshal(msg[4], &v5); err != nil {
		return err
	}

	t.V1 = v1
	t.V2 = v2
	t.V3 = v3
	t.V4 = v4
	t.V5 = v5

	return nil
}
