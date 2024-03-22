package handlers

import (
	"context"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/muneebhashone/go-fiber-api/config"
	"github.com/muneebhashone/go-fiber-api/db"
	"github.com/muneebhashone/go-fiber-api/types"
	"github.com/muneebhashone/go-fiber-api/utils"
)

var validate *validator.Validate

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(store db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: store,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userStore.GetUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{"message": "User has been deleted"})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		input types.UpdateUserInput
		id    = c.Params("id")
	)
	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := c.BodyParser(&input); err != nil {
		return err
	}

	err := validate.Struct(input)
	if err != nil {
		return err
	}

	err = h.userStore.UpdateUser(c.Context(), id, input)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{"message": "User has been updated"})
}

func (h *UserHandler) HandleLogin(c *fiber.Ctx) error {
	// Parse and validate the request body
	var input types.LoginInput

	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validate the input
	err := validate.Struct(input)
	if err != nil {
		return err
	}
	// (Assuming you have a validation setup. If not, you should manually check `input.Email` and `input.Password`)

	// Connect and query the MongoDB database to find the user
	// Assuming `collection` is your MongoDB collection where users are stored
	user, err := h.userStore.GetUser(c.Context(), input.Email)
	if err != nil {
		// Handle not found user, or other db errors
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if the provided password is correct
	// (Assuming you are storing hashed passwords)
	if !utils.CheckPasswordHash(input.Password, user.EncryptedPassword) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex()) // Assuming user ID is of type `primitive.ObjectID`
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	sess, err := config.Store.Get(c)
	if err != nil {
		return err
	}

	sess.Set("user_id", user.ID.Hex())
	err = sess.Save()
	if err != nil {
		return err
	}
	// Send the token as a response
	return c.JSON(fiber.Map{"token": token})
}

func (h *UserHandler) HandleLogout(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Destroy the session
	err = sess.Destroy()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	ctx := context.Background()

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	sortField := c.Query("sort_field", "name") // Default sorting field
	sortOrder := c.Query("sort_order", "asc")  // Default sorting order
	searchQuery := c.Query("search_query", "") // Search query

	users, total, err := h.userStore.GetUsers(ctx, page, pageSize, sortField, sortOrder, searchQuery)
	if err != nil {
		return err
	}

	paginator := map[string]interface{}{
		"current_page": page,
		"has_next":     page*pageSize < int(total),
		"total":        total,
		"page_size":    pageSize,
	}

	response := fiber.Map{
		"results":   users,
		"paginator": paginator,
	}

	return c.JSON(response)
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var input types.CreateUserInput
	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := c.BodyParser(&input); err != nil {
		return err
	}

	err := validate.Struct(input)
	if err != nil {
		return err
	}

	newUser, err := types.NewUser(input)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.CreateUser(c.Context(), *newUser)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}
