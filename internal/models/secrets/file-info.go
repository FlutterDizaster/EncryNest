package secrets

import "github.com/google/uuid"

type FileInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Description string    `json:"description,omitempty"`
}

func (f FileInfo) ConvertToSecret() (*Secret, error) {
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
