package tuple

import (
	"bytes"
	"encoding/gob"
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

func (t T2[V1, V2]) MarshalText() ([]byte, error) {
	return json.Marshal(t)
}

func (t *T2[V1, V2]) UnmarshalText(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *T2[V1, V2]) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(t.V1); err != nil {
		return []byte{}, err
	}

	if err := enc.Encode(t.V2); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (t *T2[V1, V2]) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	var v1 V1
	if err := dec.Decode(&v1); err != nil {
		return err
	}

	var v2 V2
	if err := dec.Decode(&v2); err != nil {
		return err
	}

	t.V1 = v1
	t.V2 = v2

	return nil
}

func (t *T2[V1, V2]) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

func (t *T2[V1, V2]) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}
