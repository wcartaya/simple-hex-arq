package domain

type UserService interface {
	Find(userId string) (*User, error)
	Store(user *User) error
}
