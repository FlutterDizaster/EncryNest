package secrets

// import "github.com/google/uuid"

// type CreditCard struct {
// 	ID              uuid.UUID `json:"id"`
// 	Number          string    `json:"number"`
// 	ExpirationMonth int       `json:"expiration_month"`
// 	ExpirationYear  int       `json:"expiration_year"`
// 	CVV             int       `json:"cvv"`
// 	Name            string    `json:"name"`
// 	Description     string    `json:"description,omitempty"`
// }

// func (c CreditCard) ConvertToSecret() (*Secret, error) {
// 	s := &Secret{
// 		ID:   c.ID,
// 		Kind: SecretKindCreditCard,
// 	}

// 	data, err := parseFromExactType(c)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s.Data = data

// 	return s, nil
// }
