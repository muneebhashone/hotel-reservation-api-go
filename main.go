package main

import (
	"context"
	"flag"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/muneebhashone/go-fiber-api/config"
	"github.com/muneebhashone/go-fiber-api/db"
	"github.com/muneebhashone/go-fiber-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi = "mongodb://localhost:27017"
)

func main() {
	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey:    config.JWT_SECRET,
		SigningMethod: "HS256",
		ErrorHandler: fiber.ErrorHandler(func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}),
	})

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	UserHandler := handlers.NewUserHandler(db.NewMongoUserStore(client))

	port := flag.String("port", ":5000", "Port address")

	flag.Parse()

	app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	}})

	v1 := app.Group("/api/v1")

	v1.Get("/users/:id", jwtMiddleware, UserHandler.HandleGetUser)
	v1.Post("/login", UserHandler.HandleLogin)
	v1.Delete("/users/:id", UserHandler.HandleDeleteUser)
	v1.Put("/users/:id", UserHandler.HandleUpdateUser)
	v1.Get("/users", jwtMiddleware, UserHandler.HandleGetUsers)
	v1.Post("/users", UserHandler.HandleCreateUser)

	app.Listen(*port)
}
