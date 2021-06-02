package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/prastamaha/auth-basic/database"
	"github.com/prastamaha/auth-basic/routes"
)

func main() {
	// connect into database
	creds := &database.Creds{
		Username: database.Username,
		Password: database.Password,
		Database: database.Database,
		Address:  database.Address,
		Port:     database.Port,
	}
	if err := database.Connect(creds); err != nil {
		panic(err)
	}

	// create new fiber app
	app := fiber.New()

	// add logging middleware
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
		// Output: ,
	}))

	// allows different backend and frontend
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	// setup routes
	routes.Setup(app)

	// run fiber app
	app.Listen(":3000")
}
