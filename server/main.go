package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dedsecrattle/todo-application/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Todos []models.Todo = []models.Todo{}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " +
			"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

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
				Todos[i].Done = !Todos[i].Done
				break
			}
		}

		return c.JSON(Todos)
	})

	app.Listen(":4000")
}

func CreateTodo(title, body string, todoCollection mongo.Collection) (*models.Todo, error) {
	newTodo := models.Todo{
		Title: title,
		Body:  body,
		Done:  false,
	}
	result, err := todoCollection.InsertOne(context.TODO(), newTodo)
	if err != nil {
		return nil, err
	}
	newTodo.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return &newTodo, nil
}
