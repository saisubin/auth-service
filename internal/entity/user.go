// internal/entity/user.go
package entity

import "golang.org/x/crypto/bcrypt"

type User struct {
    ID       int64
    Email    string
    Password string
}

func (u *User) HashPassword() error {
    bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
    if err != nil {
        return err
    }
    u.Password = string(bytes)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}