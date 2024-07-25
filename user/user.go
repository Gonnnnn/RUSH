package user

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	University string `json:"university"`
	Phone      string `json:"phone"`
	Generation string `json:"generation"`
	IsActive   bool   `json:"is_active"`
}
