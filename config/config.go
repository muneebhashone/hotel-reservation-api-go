package config

import "github.com/gofiber/fiber/v2/middleware/session"

const (
	DBNAME          = "hotel-reservation"
	USER_COLLECTION = "users"
)

var JWT_SECRET = []byte("your_strong_secret_key")

var Store = session.New()
