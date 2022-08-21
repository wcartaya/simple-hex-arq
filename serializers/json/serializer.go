package json

import (
	"encoding/json"
	"fmt"
	"users-ms/domain"
)

type User struct{}

func (r *User) Decode(input []byte) (*domain.User, error) {
	decoded := &domain.User{}
	if err := json.Unmarshal(input, decoded); err != nil {
		return nil, fmt.Errorf("serializer.Decode %w", err)
	}
	return decoded, nil
}

func (r *User) Encode(input *domain.User) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("serializer.Encode %w", err)
	}
	return rawMsg, nil
}
