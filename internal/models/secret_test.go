package models

// import (
// 	"bytes"
// 	"encoding/gob"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func Test_parseToExactType(t *testing.T) {
// 	type test struct {
// 		name string
// 		want interface{}
// 	}
// 	tests := []test{
// 		{
// 			name: "Credentials",
// 			want: Credentials{
// 				ID:          uuid.New(),
// 				WebSite:     "www.google.com",
// 				Username:    "username",
// 				Password:    "password",
// 				Description: "some description",
// 			},
// 		},
// 		{
// 			name: "CreditCard",
// 			want: CreditCard{
// 				ID:              uuid.New(),
// 				Number:          "1234 5678 9012 3456",
// 				ExpirationMonth: 12,
// 				ExpirationYear:  2025,
// 				CVV:             123,
// 				Name:            "name",
// 				Description:     "some description",
// 			},
// 		},
// 		{
// 			name: "Memo",
// 			want: Memo{
// 				ID:          uuid.New(),
// 				Name:        "name",
// 				Data:        "data",
// 				Description: "some description",
// 			},
// 		},
// 		{
// 			name: "FileInfo",
// 			want: FileInfo{
// 				ID:          uuid.New(),
// 				Name:        "name",
// 				Size:        12345,
// 				Description: "some description",
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			buf := &bytes.Buffer{}
// 			encoder := gob.NewEncoder(buf)

// 			switch s := tt.want.(type) {
// 			case Credentials:
// 				err := encoder.Encode(s)

// 				require.NoError(t, err, "encoder.Encode() error = %v", err)

// 				got, err := parseToExactType[Credentials](buf.Bytes())

// 				require.NoError(t, err, "parseToExactType() error = %v", err)

// 				assert.Equal(t, s, got)

// 			case CreditCard:
// 				err := encoder.Encode(s)

// 				require.NoError(t, err, "encoder.Encode() error = %v", err)

// 				got, err := parseToExactType[CreditCard](buf.Bytes())

// 				require.NoError(t, err, "parseToExactType() error = %v", err)

// 				assert.Equal(t, s, got)

// 			case Memo:
// 				err := encoder.Encode(s)

// 				require.NoError(t, err, "encoder.Encode() error = %v", err)

// 				got, err := parseToExactType[Memo](buf.Bytes())

// 				require.NoError(t, err, "parseToExactType() error = %v", err)

// 				assert.Equal(t, s, got)

// 			case FileInfo:
// 				err := encoder.Encode(s)

// 				require.NoError(t, err, "encoder.Encode() error = %v", err)

// 				got, err := parseToExactType[FileInfo](buf.Bytes())

// 				require.NoError(t, err, "parseToExactType() error = %v", err)

// 				assert.Equal(t, s, got)
// 			}
// 		})
// 	}
// }
