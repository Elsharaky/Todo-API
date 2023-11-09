package controllers

import (
	"TodoAPI/models"
	"TodoAPI/utilities"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticationMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Message": "Access token is required",
		})
	}
	parsedToken := utilities.ParseAccessToken(token)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"Expired": true,
		})
	}
	c.Locals("email", claims["email"])
	return c.Next()
}

func Register(c *fiber.Ctx) error {
	var (
		userRequest, userData models.User
		err                   error
	)

	if err = c.BodyParser(&userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	err = models.GetUser(userRequest.Email, &userData)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	if userData.Email == userRequest.Email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "Email already exists.",
		})
	}

	userData = userRequest

	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 16)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	userData.Password = string(hash)

	if err = models.InsertUser(&userData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "User registered successfully.",
	})
}

func Login(c *fiber.Ctx) error {
	var (
		userRequest, userData models.User
		err                   error
	)

	// Parse user request.
	if err = c.BodyParser(&userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	// Get user data from database and validate it.
	err = models.GetUser(userRequest.Email, &userData)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	} else if errors.Is(err, mongo.ErrNoDocuments) || bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(userRequest.Password)) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "Invalid Email or Password.",
		})
	}

	// Generate JWT access token.
	token, err := utilities.CreateAccessToken(userData.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Logged in successfully.",
		"Token":   token,
		"User":    userData,
	})
}

func Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Logged out successfully.",
		"Token":   "",
	})
}
