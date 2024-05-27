package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	//coll := client.Database("todo-app").Collection("todos")

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello this is Working!")
	})

	app.Post("/todo", func(c *fiber.Ctx) error {
		var todo models.Todo
		c.BodyParser(&todo)
		result, err := CreateTodo(&todo, client)

		if err != nil {
			fmt.Print("Error inserting Todo ", err)
		}
		return c.JSON(result)
	})

	app.Get("/todo", func(c *fiber.Ctx) error {
		result, err := GetTodos(client)

		if err != nil {
			fmt.Print("Unable to fetch Todos")
		}
		return c.JSON(result)
	})

	app.Patch("/todo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err != nil {
			return c.Status(401).SendString("Invalid Id")
		}
		UpdateTodoById(id, client)
		result, err := GetTodos(client)

		if err != nil {
			fmt.Print("Unable to fetch Todos")
		}
		return c.JSON(result)
	})

	app.Delete("/todo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err != nil {
			return c.Status(401).SendString("Invalid Id")
		}
		DeleteTodoById(id, client)
		result, err := GetTodos(client)

		if err != nil {
			fmt.Print("Unable to fetch Todos")
		}
		return c.JSON(result)
	})

	app.Listen(":4000")
}

func CreateTodo(newTodo *models.Todo, client *mongo.Client) (*models.Todo, error) {
	collection := client.Database("todo-app").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, newTodo)
	if err != nil {
		return nil, err
	}
	insertedID := result.InsertedID.(primitive.ObjectID)
	newTodo.ID = insertedID
	return newTodo, nil
}

func GetTodos(client *mongo.Client) ([]models.Todo, error) {
	collection := client.Database("todo-app").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []models.Todo
	for cursor.Next(ctx) {
		var todo models.Todo
		err := cursor.Decode(&todo)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func UpdateTodoById(id string, client *mongo.Client) error {
	collection := client.Database("todo-app").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	queryId, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": queryId}
	var todo models.Todo
	err = collection.FindOne(ctx, filter).Decode(&todo)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"done": !todo.Done}}
	defer cancel()
	if err != nil {
		fmt.Println("Error parsing the Id")
	}
	_, err = collection.UpdateOne(ctx, filter, update)

	if err != nil {
		fmt.Println("Unable to update the Todo")
	}

	return nil
}

func DeleteTodoById(id string, client *mongo.Client) error {
	collection := client.Database("todo-app").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	queryId, err := primitive.ObjectIDFromHex(id)
	defer cancel()
	if err != nil {
		fmt.Println("Unable to Parse ID")
	}
	filter := bson.M{"_id": queryId}
	_, err = collection.DeleteOne(ctx, filter)

	return nil

}
