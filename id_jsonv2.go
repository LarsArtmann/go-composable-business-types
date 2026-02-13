//go:build goexperiment.jsonv2

package cbt

import (
	jsonv2 "encoding/json/v2"
	"encoding/json/jsontext"
	"fmt"
)

// MarshalJSONTo implements jsonv2.MarshalerTo for efficient streaming JSON encoding.
// Zero values serialize to JSON null, non-zero values serialize as JSON strings.
func (id ID[B, V]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if id.IsZero() {
		return enc.WriteToken(jsontext.Null)
	}
	return enc.WriteToken(jsontext.String(id.String()))
}

// UnmarshalJSONFrom implements jsonv2.UnmarshalerFrom for efficient streaming JSON decoding.
func (id *ID[B, V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}

	if string(val) == "null" {
		var zero V
		*id = ID[B, V]{value: zero}
		return nil
	}

	var s string
	if err := jsonv2.Unmarshal(val, &s); err != nil {
		return fmt.Errorf("id: cannot unmarshal %s into %T", string(val), *id)
	}

	var zero V
	switch any(zero).(type) {
	case string:
		*id = ID[B, V]{value: any(s).(V)}
		return nil
	default:
		return fmt.Errorf("id: cannot unmarshal string into %T (only string-based IDs supported)", zero)
	}
}

// Compile-time interface assertions for jsonv2
var (
	_ jsonv2.MarshalerTo     = ID[struct{}, string]{}
	_ jsonv2.UnmarshalerFrom = (*ID[struct{}, string])(nil)
)
