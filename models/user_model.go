package models

type User struct {
	ID       string `json:"ID" bson:"_id" msgpack:"_id"`
	Username string `json:"Username" bson:"Username" msgpack:"Username"`
	Email    string `json:"Email" bson:"Email" msgpack:"Email"`
	Password string `json:"Password" bson:"Password" msgpack:"Password"`
	Name     string `json:"Name" bson:"Name" msgpack:"Name"`
	Address  string `json:"Address" bson:"Address" msgpack:"Address"`
	IsActive bool   `json:"IsActive" bson:"IsActive" msgpack:"IsActive"`
}

func NewUser() *User {
	m := new(User)
	m.IsActive = true
	return m
}

func (m *User) TableName() string {
	return "users"
}

func NewUserDefaultData() *User {
	user := NewUser()
	user.ID = "u001"
	user.Username = "user001"
	user.Password = "Password.1"
	user.Name = "User 001"
	return user
}
