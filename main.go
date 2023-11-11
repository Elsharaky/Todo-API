package main

import (
	"TodoAPI/controllers"
	"TodoAPI/initializers"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	initializers.InitDB()
}

func main() {
	defer initializers.Client.Disconnect(context.TODO())

	app := fiber.New()
	app.Use(cors.New())
	app.Use("/user", controllers.AuthenticationMiddleware)
	app.Use("/todo", controllers.AuthenticationMiddleware)

	// Authentication routes
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	// User routes
	app.Get("/user/logout", controllers.Logout)
	app.Delete("/user/delete", controllers.DeleteUser)
	app.Patch("/user/update", controllers.UpdateUser)

	// Todo routes
	app.Post("/todo/create", controllers.CreateTodo)
	app.Get("/todo/:status?", controllers.GetAllTodos)
	app.Delete("/todo/delete/:id", controllers.DeleteTodo)
	app.Patch("/todo/update/:id", controllers.UpdateTodo)

	app.Listen(":8080")

}
