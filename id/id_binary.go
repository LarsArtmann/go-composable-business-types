package id

import (
	"encoding"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

// Byte sizes for binary marshaling of integer types
const (
	byteSizeInt16 = 2 // size of int16 and uint16 in bytes
	byteSizeInt32 = 4 // size of int32 and uint32 in bytes
	byteSizeInt64 = 8 // size of int64, uint, and uint64 in bytes
)

// MarshalBinary implements encoding.BinaryMarshaler for binary encoding.
func (id ID[B, V]) MarshalBinary() ([]byte, error) {
	if id.IsZero() {
		return nil, nil
	}

	switch v := any(id.value).(type) {
	case string:
		return []byte(v), nil
	case int:
		b := make([]byte, byteSizeInt64)
		//nolint:gosec // G115: int to uint64 is safe for binary serialization
		binary.LittleEndian.PutUint64(b, uint64(v))
		return b, nil
	case int8:
		return []byte{byte(v)}, nil //nolint:gosec // G115: int8 to byte is safe for serialization
	case int16:
		b := make([]byte, byteSizeInt16)
		//nolint:gosec // G115: int16 to uint16 is safe for binary serialization
		binary.LittleEndian.PutUint16(b, uint16(v))
		return b, nil
	case int32:
		b := make([]byte, byteSizeInt32)
		//nolint:gosec // G115: int32 to uint32 is safe for binary serialization
		binary.LittleEndian.PutUint32(b, uint32(v))
		return b, nil
	case int64:
		b := make([]byte, byteSizeInt64)
		//nolint:gosec // G115: int64 to uint64 is safe for binary serialization
		binary.LittleEndian.PutUint64(b, uint64(v))
		return b, nil
	case uint:
		b := make([]byte, byteSizeInt64)
		//nolint:gosec // G115: uint to uint64 is safe for binary serialization
		binary.LittleEndian.PutUint64(b, uint64(v))
		return b, nil
	case uint8:
		return []byte{v}, nil
	case uint16:
		b := make([]byte, byteSizeInt16)
		binary.LittleEndian.PutUint16(b, v)
		return b, nil
	case uint32:
		b := make([]byte, byteSizeInt32)
		binary.LittleEndian.PutUint32(b, v)
		return b, nil
	case uint64:
		b := make([]byte, byteSizeInt64)
		binary.LittleEndian.PutUint64(b, v)
		return b, nil
	default:
		return nil, fmt.Errorf("id: unsupported type %T for binary marshaling", id.value)
	}
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for binary decoding.
func (id *ID[B, V]) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		id.Reset()
		return nil
	}

	var zero V
	switch any(zero).(type) {
	case string:
		*id = ID[B, V]{value: any(string(data)).(V)}
		return nil
	case int:
		if len(data) < byteSizeInt64 {
			return fmt.Errorf(
				"id: insufficient data for int: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt64,
				data,
				zero,
			)
		}
		n := int(
			binary.LittleEndian.Uint64(data),
		)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int8:
		if len(data) < 1 {
			return fmt.Errorf(
				"id: insufficient data for int8: got %d bytes, want 1 (data=%x, targetType=%T)",
				len(data),
				data,
				zero,
			)
		}
		*id = ID[B, V]{value: any(int8(data[0])).(V)}
		return nil
	case int16:
		if len(data) < byteSizeInt16 {
			return fmt.Errorf(
				"id: insufficient data for int16: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt16,
				data,
				zero,
			)
		}
		n := int16(
			binary.LittleEndian.Uint16(data),
		)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int32:
		if len(data) < byteSizeInt32 {
			return fmt.Errorf(
				"id: insufficient data for int32: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt32,
				data,
				zero,
			)
		}
		n := int32(
			binary.LittleEndian.Uint32(data),
		)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int64:
		if len(data) < byteSizeInt64 {
			return fmt.Errorf(
				"id: insufficient data for int64: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt64,
				data,
				zero,
			)
		}
		n := int64(
			binary.LittleEndian.Uint64(data),
		)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint:
		if len(data) < byteSizeInt64 {
			return fmt.Errorf(
				"id: insufficient data for uint: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt64,
				data,
				zero,
			)
		}
		n := uint(binary.LittleEndian.Uint64(data))
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint8:
		if len(data) < 1 {
			return fmt.Errorf(
				"id: insufficient data for uint8: got %d bytes, want 1 (data=%x, targetType=%T)",
				len(data),
				data,
				zero,
			)
		}
		*id = ID[B, V]{value: any(data[0]).(V)}
		return nil
	case uint16:
		if len(data) < byteSizeInt16 {
			return fmt.Errorf(
				"id: insufficient data for uint16: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt16,
				data,
				zero,
			)
		}
		n := binary.LittleEndian.Uint16(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint32:
		if len(data) < byteSizeInt32 {
			return fmt.Errorf(
				"id: insufficient data for uint32: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt32,
				data,
				zero,
			)
		}
		n := binary.LittleEndian.Uint32(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint64:
		if len(data) < byteSizeInt64 {
			return fmt.Errorf(
				"id: insufficient data for uint64: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt64,
				data,
				zero,
			)
		}
		n := binary.LittleEndian.Uint64(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	default:
		return fmt.Errorf("id: unsupported type %T for binary unmarshaling (data=%x)", zero, data)
	}
}

// GobEncode implements gob.GobEncoder for Go-specific encoding.
func (id ID[B, V]) GobEncode() ([]byte, error) {
	return id.MarshalBinary()
}

// GobDecode implements gob.GobDecoder for Go-specific decoding.
func (id *ID[B, V]) GobDecode(data []byte) error {
	return id.UnmarshalBinary(data)
}

// Compile-time interface assertions for binary encoding
var (
	_ encoding.BinaryMarshaler   = ID[struct{}, string]{value: ""}
	_ encoding.BinaryUnmarshaler = (*ID[struct{}, string])(nil)
	_ encoding.BinaryMarshaler   = ID[struct{}, int64]{value: 0}
	_ encoding.BinaryUnmarshaler = (*ID[struct{}, int64])(nil)
	_ gob.GobEncoder             = ID[struct{}, string]{value: ""}
	_ gob.GobDecoder             = (*ID[struct{}, string])(nil)
)
