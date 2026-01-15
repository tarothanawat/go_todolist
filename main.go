package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// 1 hr 03:00
type Todo struct {
	// MongoDB uses _id as the primary key field, and has type ObjectID == bson.ObjectID
	ID bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool `json:"completed"`
	Body string `json:"body"`

}

var collection *mongo.Collection
func main() {
	fmt.Println("Hello world")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	MONGODB_URI := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().ApplyURI(MONGODB_URI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
  client, err := mongo.Connect(clientOptions)
  if err != nil {
    panic(err)
  }

	defer func() {
    if err = client.Disconnect(context.Background()); err != nil {
      panic(err)
    }
  }()
  // Send a ping to confirm a successful connection
  if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
    panic(err)
  }
  fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	collection = client.Database("golang_db").Collection("todos")
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen(":" + port))

}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	// Create a context with timeout
	// Timer that says giveup after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
	// defer is guranteed to run at the end of the function
	defer cancel()
	// query in mongoDB returns a cursor object which is a pointer to the result set of a query
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into a Todo struct
	for cursor.Next(ctx) {
		var todo Todo
		// &todo is the address of the todo variable, means put the decoded data into this variable
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	// Check for errors during iteration
	if err := cursor.Err(); err != nil {
		return err
	}

	return c.JSON(todos)
}
func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	// binds a request body to a struct
	if err := c.BodyParser(todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Body is required"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(bson.ObjectID)
	// 201 response code means created
	return c.Status(201).JSON(todo)
}
func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	var updatedTodo Todo
	err = collection.FindOne(context.Background(), filter).Decode(&updatedTodo)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(updatedTodo)
}
func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}
	
	_, err = collection.DeleteOne(context.Background(), filter)

	return c.Status(200).JSON(fiber.Map{"message": "Todo deleted successfully"})
}