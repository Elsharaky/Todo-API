package controllers

import (
	"TodoAPI/models"

	"github.com/gofiber/fiber/v2"
)

func CreateTodo(c *fiber.Ctx) error {
	var (
		todo models.Todo
	)
	email := c.Locals("email").(string)

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	if err := models.CreateTodo(email, todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Todo created successfully.",
	})
}

func GetAllTodos(c *fiber.Ctx) error {
	status := c.Params("status")
	email := c.Locals("email").(string)

	var todos []models.Todo

	if err := models.GetTodos(email, status, &todos); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Todos": todos,
	})
}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	email := c.Locals("email").(string)

	if err := models.DeleteTodo(email, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Todo deleted successfully.",
	})
}

func UpdateTodo(c *fiber.Ctx) error {
	var todo models.Todo

	email := c.Locals("email").(string)
	id := c.Params("id")

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	if err := models.UpdateTodo(email, id, todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Todo updated successfully.",
	})
}
