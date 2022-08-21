package domain

type UserSerializer interface {
	Decode(input []byte) (*User, error)
	Encode(input *User) ([]byte, error)
}
