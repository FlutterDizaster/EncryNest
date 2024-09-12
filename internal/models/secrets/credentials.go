package secrets

import "github.com/google/uuid"

type Credentials struct {
	ID          uuid.UUID `json:"id"`
	WebSite     string    `json:"website"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Description string    `json:"description,omitempty"`
}

func (c Credentials) ConvertToSecret() (*Secret, error) {
	s := &Secret{
		ID:   c.ID,
		Kind: SecretKindCredentials,
	}

	data, err := parseFromExactType(c)
	if err != nil {
		return nil, err
	}

	s.Data = data

	return s, nil
}
