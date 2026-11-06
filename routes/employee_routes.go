package routes

import (
	"crud-employee/controllers"

	"github.com/gin-gonic/gin"
)

// SetupEmployeeRoutes registers all employee API endpoints with the Gin router.
//
// GOLANG CONCEPT - HOW THE GIN FRAMEWORK WORKS:
// Gin is a high-performance HTTP web framework written in Go.
// It manages HTTP routing, request parsing, middleware handling, and JSON responses.
//
// GOLANG CONCEPT - HTTP METHOD ROUTING:
// In REST APIs, different HTTP verbs correspond to CRUD actions:
// - router.GET(): Used to retrieve data (Read)
// - router.POST(): Used to create new data (Create)
// - router.PUT(): Used to update existing data (Update)
// - router.DELETE(): Used to remove data (Delete)
//
// Gin matches incoming HTTP requests by both method and URL pattern,
// then invokes the assigned controller handler function.
func SetupEmployeeRoutes(router *gin.Engine) {
	// Create a route group with prefix "/api"
	// Grouping helps organize routes and apply shared middleware or versioning.
	api := router.Group("/api")
	{
		// GET /api/employees - Retrieve all employees
		api.GET("/employees", controllers.GetAllEmployees)

		// GET /api/employees/:id - Retrieve a single employee by ID
		// The ":id" syntax defines a named path parameter extracted via c.Param("id")
		api.GET("/employees/:id", controllers.GetEmployeeByID)

		// POST /api/employees - Create a new employee
		api.POST("/employees", controllers.CreateEmployee)

		// PUT /api/employees/:id - Update an employee by ID
		api.PUT("/employees/:id", controllers.UpdateEmployee)

		// DELETE /api/employees/:id - Delete an employee by ID
		api.DELETE("/employees/:id", controllers.DeleteEmployee)
	}
}
