package models

import (
	"errors"

	pb "github.com/FlutterDizaster/EncryNest/api/generated"
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
	SecretKindFileInfo
)

// type ISecret interface {
// 	ConvertToSecret() (*Secret, error)
// }

type Secret struct {
	Version string
	ID      uuid.UUID  `json:"id"`
	Kind    SecretKind `json:"kind"`
	Data    []byte     `json:"data"`
}

func NewSecretFromProto(secret *pb.Secret) (*Secret, error) {
	result := &Secret{}

	id, err := uuid.Parse(secret.GetId())
	if err != nil {
		return nil, err
	}
	result.ID = id
	result.Data = secret.GetData()
	result.Version = secret.GetVersion()

	// Determining secret kind
	switch secret.GetKind() {
	case pb.SecretKind_CREDENTIALS:
		result.Kind = SecretKindCredentials
	case pb.SecretKind_CREDIT_CARD:
		result.Kind = SecretKindCreditCard
	case pb.SecretKind_MEMO:
		result.Kind = SecretKindMemo
	case pb.SecretKind_FILE_INFO:
		result.Kind = SecretKindFileInfo
	}

	return result, nil
}

func (s *Secret) ToProto() *pb.Secret {
	secret := &pb.Secret{}

	secret.Id = s.ID.String()
	secret.Data = s.Data
	secret.Version = s.Version

	switch s.Kind {
	case SecretKindCredentials:
		secret.Kind = pb.SecretKind_CREDENTIALS
	case SecretKindCreditCard:
		secret.Kind = pb.SecretKind_CREDIT_CARD
	case SecretKindMemo:
		secret.Kind = pb.SecretKind_MEMO
	case SecretKindFileInfo:
		secret.Kind = pb.SecretKind_FILE_INFO
	}

	return secret
}

// // ParseToExactType parses secret to exact type.
// func (s *Secret) ParseToExactType() (interface{}, error) {
// 	switch s.Kind {
// 	case SecretKindCredentials:
// 		return parseToExactType[Credentials](s.Data)

// 	case SecretKindCreditCard:
// 		return parseToExactType[CreditCard](s.Data)

// 	case SecretKindMemo:
// 		return parseToExactType[Memo](s.Data)

// 	case SecretKindFileInfo:
// 		return parseToExactType[FileInfo](s.Data)

// 	default:
// 		return nil, ErrWrongType
// 	}
// }

// // parseToExactType decodes data into a value of type T.
// func parseToExactType[T ISecret](data []byte) (T, error) {
// 	var result T
// 	buf := bytes.NewBuffer(data)
// 	decoder := gob.NewDecoder(buf)
// 	err := decoder.Decode(&result)
// 	if err != nil {
// 		return result, err
// 	}
// 	return result, nil
// }

// func parseFromExactType[T ISecret](data T) ([]byte, error) {
// 	var buf bytes.Buffer
// 	encoder := gob.NewEncoder(&buf)
// 	err := encoder.Encode(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }
