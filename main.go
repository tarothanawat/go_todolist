package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	// the backticks are the actual thing you'll actually see in the json
	ID int `json:"id"`
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

func main() {
	fmt.Println("Hello worldx")
	app := fiber.New()

	// pointers
		// var x int = 5 // 0x000001
		// var y *int = &x // y is pointer to an int, &x is the address of x

		// fmt.Println(y) // 0x000001
		// fmt.Println(*y) // 5, dereferencing the pointer
	// in-memory todos slice
	todos := []Todo{}

	// GET /
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{} // {id:0, completed:false, body:""}

		if error := c.BodyParser(todo); error != nil {
			return error
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body is required"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	// Update a Todo to Completed
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})	
	})

	// Delete a Todo

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}
		// 1 2 3 4 5
		// 1 2 4 5
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})
	log.Fatal(app.Listen(":4000"))
}