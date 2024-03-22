package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_COST = 12
)

type CreateUserInput struct {
	Firstname string `json:"firstname" validate:"required,min=3,max=32"`
	Lastname  string `json:"lastname" validate:"required,min=3,max=32"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=64"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserInput struct {
	Firstname string    `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string    `json:"lastname,omitempty" bson:"lastname,omitempty"`
	UpdatedAt time.Time `json:"-" bson:"updated_at"`
}

type User struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Firstname         string             `json:"firstname" bson:"firstname"`
	Lastname          string             `json:"lastname" bson:"lastname"`
	Email             string             `json:"email" bson:"email"`
	EncryptedPassword string             `json:"-" bson:"encrypted_password"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewUser(input CreateUserInput) (*User, error) {
	currentTime := time.Now()

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), BCRYPT_COST)
	if err != nil {
		return nil, err
	}

	return &User{
		Firstname:         input.Firstname,
		Lastname:          input.Lastname,
		Email:             input.Email,
		EncryptedPassword: string(encryptedPassword),
		CreatedAt:         currentTime,
		UpdatedAt:         currentTime,
	}, nil
}
