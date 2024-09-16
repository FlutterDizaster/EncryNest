package secrets

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/google/uuid"
)

type SecretKind int

// Errors.
var (
	ErrWrongType = errors.New("wrong type")
)

// Secrets kinds.
const (
	SecretKindCredentials SecretKind = iota
	SecretKindCreditCard
	SecretKindMemo
	SecretKindFilePart
)

type ISecret interface {
	ConvertToSecret() (*Secret, error)
}

type Secret struct {
	ID   uuid.UUID  `json:"id"`
	Kind SecretKind `json:"kind"`
	Data []byte     `json:"data"`
}

// ParseToExactType parses secret to exact type.
func (s *Secret) ParseToExactType() (interface{}, error) {
	switch s.Kind {
	case SecretKindCredentials:
		return parseToExactType[Credentials](s.Data)

	case SecretKindCreditCard:
		return parseToExactType[CreditCard](s.Data)

	case SecretKindMemo:
		return parseToExactType[Memo](s.Data)

	case SecretKindFilePart:
		return parseToExactType[FileInfo](s.Data)

	default:
		return nil, ErrWrongType
	}
}

// parseToExactType decodes data into a value of type T.
func parseToExactType[T ISecret](data []byte) (T, error) {
	var result T
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func parseFromExactType[T ISecret](data T) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
