package main

import (
	"github.com/dedsecrattle/todo-application/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var Todos []models.Todo = []models.Todo{}

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello this is Working!")
	})

	app.Post("/todo", func(c *fiber.Ctx) error {
		var todo models.Todo
		c.BodyParser(&todo)
		todo.ID = len(Todos) + 1
		Todos = append(Todos, todo)
		return c.JSON(Todos)
	})

	app.Get("/todo", func(c *fiber.Ctx) error {
		return c.JSON(Todos)
	})

	app.Patch("/todo/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(401).SendString("Invalid Id")
		}

		for i, t := range Todos {
			if t.ID == id {
				Todos[i].Done = true
				break
			}
		}

		return c.JSON(Todos)

	})
	app.Listen(":4000")
}
