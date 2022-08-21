package domain

type User struct {
	UserId    string `json:"code" bson:"code" msgpack:"code"`
	Name      string `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at"`
}
