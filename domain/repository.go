package domain

type UserRepository interface {
	Find(userId string) (*User, error)
	Store(user *User) error
}
