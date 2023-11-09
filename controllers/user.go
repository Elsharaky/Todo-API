package controllers

import (
	"TodoAPI/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func DeleteUser(c *fiber.Ctx) error {
	var (
		userData models.User
		err      error
	)
	email := c.Locals("email").(string)

	// Make sure that the user already exists.
	err = models.GetUser(email, &userData)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "User does not exist.",
		})
	}

	if err = models.DeleteUser(email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "User deleted successfully.",
		"Token":   "",
	})
}

func UpdateUser(c *fiber.Ctx) error {
	var (
		userRequest models.User
	)

	email := c.Locals("email").(string)
	if err := c.BodyParser(&userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	if userRequest.Email != "" && userRequest.Email != email {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "You cannot edit email address.",
		})
	}

	if userRequest.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 16)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Error": err.Error(),
			})
		}
		userRequest.Password = string(hash)
	}

	userRequest.Todos = []models.Todo{}
	if err := models.UpdateUser(email, userRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "User updated successfully.",
	})
}
