package models

type User struct {
	ID       string `json:"ID" bson:"_id" msgpack:"_id" db:"ID"`
	Name     string `json:"Name" bson:"Name" msgpack:"Name" db:"Name"`
	Username string `json:"Username" bson:"Username" msgpack:"Username" db:"Username"`
	Email    string `json:"Email" bson:"Email" msgpack:"Email" db:"Email"`
	Password string `json:"Password" bson:"Password" msgpack:"Password" db:"Password"`
	Address  string `json:"Address" bson:"Address" msgpack:"Address" db:"Address"`
	IsActive bool   `json:"IsActive" bson:"IsActive" msgpack:"IsActive" db:"IsActive"`
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

func (user *User) GetMapFormat() map[string]interface{} {
	return map[string]interface{}{
		"ID":       user.ID,
		"Name":     user.Name,
		"Username": user.Username,
		"Email":    user.Email,
		"Password": user.Password,
		"Address":  user.Address,
		"IsActive": user.IsActive,
	}
}

func (user *User) SplitByField() []interface{} {
	return []interface{}{
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Address,
		user.IsActive,
	}
}
