package tuple

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func NewT3[V1, V2, V3 any](v1 V1, v2 V2, v3 V3) T3[V1, V2, V3] {
	return T3[V1, V2, V3]{v1, v2, v3}
}

type T3[V1, V2, V3 any] struct {
	V1 V1
	V2 V2
	V3 V3
}

func (t T3[V1, V2, V3]) Values() (V1, V2, V3) {
	return t.V1, t.V2, t.V3
}

func (t T3[V1, V2, V3]) MarshalJSON() ([]byte, error) {
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

	return json.Marshal([]json.RawMessage{b1, b2, b3})
}

func (t *T3[V1, V2, V3]) UnmarshalJSON(p []byte) error {
	var msg [3]json.RawMessage
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

	t.V1 = v1
	t.V2 = v2
	t.V3 = v3

	return nil
}

func (t T3[V1, V2, V3]) MarshalText() ([]byte, error) {
	return json.Marshal(t)
}

func (t *T3[V1, V2, V3]) UnmarshalText(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *T3[V1, V2, V3]) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(t.V1); err != nil {
		return []byte{}, err
	}

	if err := enc.Encode(t.V2); err != nil {
		return []byte{}, err
	}

	if err := enc.Encode(t.V3); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (t *T3[V1, V2, V3]) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	var v1 V1
	if err := dec.Decode(&v1); err != nil {
		return err
	}

	var v2 V2
	if err := dec.Decode(&v2); err != nil {
		return err
	}

	var v3 V3
	if err := dec.Decode(&v3); err != nil {
		return err
	}

	t.V1 = v1
	t.V2 = v2
	t.V3 = v3

	return nil
}

func (t *T3[V1, V2, V3]) GobEncode() ([]byte, error) {
	return t.MarshalBinary()
}

func (t *T3[V1, V2, V3]) GobDecode(data []byte) error {
	return t.UnmarshalBinary(data)
}
