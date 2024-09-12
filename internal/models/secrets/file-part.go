package secrets

import "github.com/google/uuid"

type FilePart struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Data        []byte    `json:"data"`
	Description string    `json:"description,omitempty"`
}

func (f FilePart) ConvertToSecret() (*Secret, error) {
	s := &Secret{
		ID:   f.ID,
		Kind: SecretKindFilePart,
	}

	data, err := parseFromExactType(f)
	if err != nil {
		return nil, err
	}

	s.Data = data

	return s, nil
}
