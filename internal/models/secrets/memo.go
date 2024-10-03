package secrets

// import "github.com/google/uuid"

// type Memo struct {
// 	ID          uuid.UUID `json:"id"`
// 	Name        string    `json:"name"`
// 	Data        string    `json:"data"`
// 	Description string    `json:"description,omitempty"`
// }

// func (m Memo) ConvertToSecret() (*Secret, error) {
// 	s := &Secret{
// 		ID:   m.ID,
// 		Kind: SecretKindMemo,
// 	}

// 	data, err := parseFromExactType(m)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s.Data = data

// 	return s, nil
// }
