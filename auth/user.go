package auth

import (
	"crypto/sha256"
	"time"

	"github.com/google/uuid"
)

// User representa o usuário do sistema.
type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUser é o construtor para a struct User.
func NewUser(username, password string) *User {
	return &User{
		Id:        uuid.New(),
		Username:  username,
		Password:  HashPassword(password),
		Timestamp: time.Now(),
	}
}

// HashPassword recebe um password e retorna a sua versão
// criptografada.
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return string(hash[:])
}
