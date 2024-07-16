package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method:
	route.Get("/getIOC:id", controllers.GetIOCs)
	route.Get("/test", test) // IoC information by IP address.
}

func test(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"data":  "sa",
	})

}
