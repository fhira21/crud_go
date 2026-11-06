package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"crud-employee/config"
	"crud-employee/models"

	"github.com/gin-gonic/gin"
)

// GOLANG CONCEPT - ERROR HANDLING WITH "if err != nil":
// Go does not use exceptions (try/catch). Instead, functions that can fail return an error
// as their last return value. By convention, if a function succeeds, error is `nil`.
// If something goes wrong, error is non-nil containing the error details.
// Developers explicitly check `if err != nil` to handle issues gracefully.

// GetAllEmployees retrieves all employee records from MySQL.
// Handled by: GET /api/employees
func GetAllEmployees(c *gin.Context) {
	// GOLANG CONCEPT - db.Query():
	// Use Query() for SQL statements that return multiple rows (typically SELECT queries).
	// It returns a *sql.Rows cursor that allows us to iterate through each record.
	query := "SELECT id, name, npwp, address FROM employee ORDER BY id DESC"
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch employees from database",
			"error":   err.Error(),
		})
		return
	}
	// Always close rows cursor when the function finishes to prevent memory leaks and free the database connection.
	defer rows.Close()

	// Initialize an empty slice of employees.
	// Using make([]models.Employee, 0) guarantees the JSON serializer outputs [] instead of null when empty.
	employees := make([]models.Employee, 0)

	// Iterate over each row returned by MySQL
	for rows.Next() {
		var emp models.Employee
		// rows.Scan copies column values from the current row into the destination variables.
		err := rows.Scan(&emp.ID, &emp.Name, &emp.NPWP, &emp.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to read employee record",
				"error":   err.Error(),
			})
			return
		}
		employees = append(employees, emp)
	}

	// Check if any error occurred during the row iteration
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error iterating employee rows",
			"error":   err.Error(),
		})
		return
	}

	// Return successful response with employee list
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Employees retrieved successfully",
		"data":    employees,
	})
}

// GetEmployeeByID retrieves a single employee by their ID.
// Handled by: GET /api/employees/:id
func GetEmployeeByID(c *gin.Context) {
	// GOLANG CONCEPT - RETRIEVING URL PARAMETERS:
	// c.Param("id") retrieves the dynamic parameter value from the URL path pattern "/api/employees/:id".
	// The returned value is always a string, so we use strconv.Atoi() to convert it to an integer.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid employee ID. Must be a positive integer.",
		})
		return
	}

	// GOLANG CONCEPT - db.QueryRow():
	// Use QueryRow() when a query is expected to return at most one single row.
	// We use parameterized queries ("?") to prevent SQL injection vulnerabilities.
	query := "SELECT id, name, npwp, address FROM employee WHERE id = ?"
	var emp models.Employee

	// QueryRow scans directly into the struct fields.
	err = config.DB.QueryRow(query, id).Scan(&emp.ID, &emp.Name, &emp.NPWP, &emp.Address)
	if err != nil {
		// If no row was found matching the ID, sql.ErrNoRows is returned.
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Employee not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database query error",
			"error":   err.Error(),
		})
		return
	}

	// Return the found employee
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Employee retrieved successfully",
		"data":    emp,
	})
}

// CreateEmployee creates a new employee record in MySQL.
// Handled by: POST /api/employees
func CreateEmployee(c *gin.Context) {
	var input models.Employee

	// GOLANG CONCEPT - ShouldBindJSON():
	// ShouldBindJSON reads the HTTP request body (JSON format) and maps its key-value pairs
	// into the fields of our Go struct according to the `json:"..."` struct tags.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	// Basic Input Validation
	input.Name = strings.TrimSpace(input.Name)
	input.NPWP = strings.TrimSpace(input.NPWP)
	input.Address = strings.TrimSpace(input.Address)

	if input.Name == "" || input.NPWP == "" || input.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validation error: Name, NPWP, and Address are all required fields.",
		})
		return
	}

	// GOLANG CONCEPT - db.Exec():
	// Use Exec() for queries that modify data and do NOT return rows, such as INSERT, UPDATE, and DELETE.
	// Exec returns a sql.Result which provides:
	// - LastInsertId() to get the auto-generated primary key ID.
	// - RowsAffected() to know how many rows were modified.
	query := "INSERT INTO employee (name, npwp, address) VALUES (?, ?, ?)"
	result, err := config.DB.Exec(query, input.Name, input.NPWP, input.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to insert employee into database",
			"error":   err.Error(),
		})
		return
	}

	// Retrieve the auto-increment ID generated by MySQL
	newID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve generated employee ID",
			"error":   err.Error(),
		})
		return
	}

	input.ID = int(newID)

	// Return 201 Created status code with the newly created employee data
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Employee created successfully",
		"data":    input,
	})
}

// UpdateEmployee updates an existing employee record.
// Handled by: PUT /api/employees/:id
func UpdateEmployee(c *gin.Context) {
	// 1. Retrieve and validate the URL path parameter ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid employee ID. Must be a positive integer.",
		})
		return
	}

	// 2. Parse the request body JSON
	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	// 3. Validate input fields
	input.Name = strings.TrimSpace(input.Name)
	input.NPWP = strings.TrimSpace(input.NPWP)
	input.Address = strings.TrimSpace(input.Address)

	if input.Name == "" || input.NPWP == "" || input.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validation error: Name, NPWP, and Address are all required fields.",
		})
		return
	}

	// 4. Verify whether the employee exists before attempting update
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM employee WHERE id = ?)"
	err = config.DB.QueryRow(checkQuery, id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Employee not found",
		})
		return
	}

	// 5. Execute the UPDATE statement with parameterized query placeholders
	updateQuery := "UPDATE employee SET name = ?, npwp = ?, address = ? WHERE id = ?"
	_, err = config.DB.Exec(updateQuery, input.Name, input.NPWP, input.Address, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update employee in database",
			"error":   err.Error(),
		})
		return
	}

	input.ID = id

	// Return 200 OK with the updated employee object
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Employee updated successfully",
		"data":    input,
	})
}

// DeleteEmployee deletes an employee record by ID.
// Handled by: DELETE /api/employees/:id
func DeleteEmployee(c *gin.Context) {
	// 1. Retrieve and validate the URL path parameter ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid employee ID. Must be a positive integer.",
		})
		return
	}

	// 2. Execute DELETE query
	deleteQuery := "DELETE FROM employee WHERE id = ?"
	result, err := config.DB.Exec(deleteQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete employee from database",
			"error":   err.Error(),
		})
		return
	}

	// 3. Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to inspect affected rows",
			"error":   err.Error(),
		})
		return
	}

	// If 0 rows were affected, the ID did not exist in the database
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Employee not found",
		})
		return
	}

	// Return 200 OK
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Employee deleted successfully",
	})
}
