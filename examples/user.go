package examples

import "github.com/9ssi7/rapidval"

type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Age       int
}

func (u *User) Validate(v *rapidval.Validator) error {
	return v.Validate(rapidval.P{
		rapidval.Required("FirstName", u.FirstName),
		rapidval.Required("LastName", u.LastName),
		rapidval.MinLength("FirstName", u.FirstName, 2),
		rapidval.MinLength("LastName", u.LastName, 2),
		rapidval.Required("Email", u.Email),
		rapidval.Email("Email", u.Email),
		rapidval.Required("Password", u.Password),
		rapidval.MinLength("Password", u.Password, 8),
		rapidval.Required("Age", u.Age),
		rapidval.Between("Age", u.Age, 18, 100),
	})
}
