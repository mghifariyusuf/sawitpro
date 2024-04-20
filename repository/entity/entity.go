package entity

type UserData struct {
	ID          int64
	PhoneNumber string
	FullName    string
	Password    string
}

type UserFilter struct {
	ID          *int64
	PhoneNumber *string
}
