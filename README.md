# Employee Management CRUD Application (Golang & MySQL)

A beginner-friendly, clean, and modern Employee Management CRUD application built with **Golang (Gin Framework)** and **MySQL** on the backend, paired with a responsive **Vanilla JavaScript (Fetch API)**, HTML5, and CSS3 frontend.

This project is structured specifically to help developers learn:
- How a standard Go project is structured.
- How to connect and interact with MySQL using Go's `database/sql` and `go-sql-driver/mysql`.
- How to develop RESTful APIs using the **Gin Web Framework**.
- How HTTP requests, JSON payloads, and responses work.
- How frontend applications communicate with Go backends using the native JavaScript `fetch()` API.

---

## 🚀 Technologies Used

- **Backend:**
  - **Golang** (Go 1.23+)
  - **Gin Web Framework** (`github.com/gin-gonic/gin`)
  - **MySQL Database**
  - **Go MySQL Driver** (`github.com/go-sql-driver/mysql`)
  - **Godotenv** (`github.com/joho/godotenv`) for environment variable management
- **Frontend:**
  - **HTML5** & **CSS3** (Modern, clean, responsive UI with light gray background and blue accent)
  - **Vanilla JavaScript** (ES6+)
  - **Fetch API** (No external frameworks or libraries)

---

## 📁 Project Structure

```text
crud-employee/
├── main.go                     # Application entry point, server setup, CORS & static files
├── go.mod                      # Go module definitions
├── go.sum                      # Checksums of Go dependencies
├── .env                        # Local database credentials and server settings
├── .env.example                # Example environment configuration template
├── .gitignore                  # Git ignore rules
│
├── config/
│   └── database.go             # MySQL connection initialization and connection pool configuration
│
├── models/
│   └── employee.go             # Employee struct and JSON serialization tags
│
├── controllers/
│   └── employee_controller.go  # CRUD handlers (GetAll, GetByID, Create, Update, Delete)
│
├── routes/
│   └── employee_routes.go      # Gin router setup and endpoint mapping
│
├── frontend/
│   ├── index.html              # Frontend user interface
│   ├── style.css               # Clean, modern UI styling
│   └── script.js               # Dynamic DOM rendering and Fetch API integration
│
└── README.md                   # Complete project documentation
```

---

## 📋 Prerequisites

Before running this project, ensure you have installed:
1. **Go** (version 1.20 or newer): [Download Go](https://go.dev/dl/)
   Verify installation:
   ```bash
   go version
   ```
2. **MySQL Server** (running locally on port `3306`):
   Verify MySQL is running:
   ```bash
   mysql -u root -p
   ```

---

## 🗄️ Database Setup

Ensure your MySQL database and table exist as specified:

- **Database Name:** `crud-employee`
- **Table Name:** `employee`

Table Schema:
```sql
CREATE DATABASE IF NOT EXISTS `crud-employee`;

USE `crud-employee`;

CREATE TABLE IF NOT EXISTS employee (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    npwp VARCHAR(255),
    address VARCHAR(255)
);
```

---

## ⚙️ Configuration (.env)

Create or update the `.env` file in the root directory:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=crud-employee
SERVER_PORT=8080
```

> **Note:** If your MySQL `root` user has a password, enter it in `DB_PASSWORD=your_password`.

---

## 📦 Installation & Setup

1. **Clone or navigate to the project directory:**
   ```bash
   cd /path/to/project
   ```

2. **Download and tidy Go dependencies:**
   ```bash
   go mod tidy
   ```

---

## ▶️ Running the Application

### 1. Start the Golang Server
Run the application with:
```bash
go run main.go
```

When the server starts successfully, you will see output like:
```text
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
Successfully connected to MySQL database: crud-employee
==================================================
🚀 Server running on http://localhost:8080
💻 Open frontend at: http://localhost:8080/
📡 API Base URL at: http://localhost:8080/api/employees
==================================================
```

### 2. Access the Frontend
Open your web browser and navigate to:
```
http://localhost:8080/
```
The Gin server serves the frontend directly from the `./frontend` directory. Alternatively, you can open `frontend/index.html` using a local web server (such as VS Code Live Server) because CORS is enabled.

---

## 📡 REST API Documentation

Base URL: `http://localhost:8080/api`

| Method | Endpoint | Description |
|:---|:---|:---|
| **GET** | `/api/employees` | Retrieve a list of all employees |
| **GET** | `/api/employees/:id` | Retrieve a single employee by their ID |
| **POST** | `/api/employees` | Create a new employee record |
| **PUT** | `/api/employees/:id` | Update an existing employee record |
| **DELETE** | `/api/employees/:id` | Delete an employee record |

---

## 🧪 Example API Requests & Responses

### 1. Retrieve All Employees
- **Request:**
  ```http
  GET /api/employees HTTP/1.1
  Host: localhost:8080
  Accept: application/json
  ```
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Employees retrieved successfully",
    "data": [
      {
        "id": 1,
        "name": "Budi Santoso",
        "npwp": "01.234.567.8-901.000",
        "address": "Jl. Sudirman No. 10, Jakarta"
      }
    ]
  }
  ```

---

### 2. Retrieve Employee by ID
- **Request:**
  ```http
  GET /api/employees/1 HTTP/1.1
  Host: localhost:8080
  ```
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Employee retrieved successfully",
    "data": {
      "id": 1,
      "name": "Budi Santoso",
      "npwp": "01.234.567.8-901.000",
      "address": "Jl. Sudirman No. 10, Jakarta"
    }
  }
  ```
- **Response if not found (404 Not Found):**
  ```json
  {
    "success": false,
    "message": "Employee not found"
  }
  ```

---

### 3. Create a New Employee
- **Request:**
  ```http
  POST /api/employees HTTP/1.1
  Host: localhost:8080
  Content-Type: application/json

  {
    "name": "Siti Nurhaliza",
    "npwp": "02.345.678.9-102.000",
    "address": "Jl. Merdeka No. 45, Bandung"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "Employee created successfully",
    "data": {
      "id": 2,
      "name": "Siti Nurhaliza",
      "npwp": "02.345.678.9-102.000",
      "address": "Jl. Merdeka No. 45, Bandung"
    }
  }
  ```

---

### 4. Update an Employee
- **Request:**
  ```http
  PUT /api/employees/2 HTTP/1.1
  Host: localhost:8080
  Content-Type: application/json

  {
    "name": "Siti Nurhaliza, S.T.",
    "npwp": "02.345.678.9-102.000",
    "address": "Jl. Diponegoro No. 88, Bandung"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Employee updated successfully",
    "data": {
      "id": 2,
      "name": "Siti Nurhaliza, S.T.",
      "npwp": "02.345.678.9-102.000",
      "address": "Jl. Diponegoro No. 88, Bandung"
    }
  }
  ```

---

### 5. Delete an Employee
- **Request:**
  ```http
  DELETE /api/employees/2 HTTP/1.1
  Host: localhost:8080
  ```
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Employee deleted successfully"
  }
  ```

---

## 💡 Golang Key Concepts Explained

### 1. `package main` & `func main()`
- Every standalone executable Go program must declare `package main`.
- `func main()` is the entry point. When you run `go run main.go`, execution starts here.

### 2. Structs & JSON Struct Tags
- A `struct` is a typed collection of fields representing a data model.
- JSON struct tags (e.g. `` `json:"npwp"` ``) dictate how struct fields are mapped to and from JSON keys.

### 3. Database Operations: `Query()`, `QueryRow()`, and `Exec()`
- **`db.Query()`**: Used for queries that return multiple records (`SELECT * FROM employee`). Returns `*sql.Rows`. Must be closed with `defer rows.Close()`.
- **`db.QueryRow()`**: Used when expecting at most one row (`SELECT ... WHERE id = ?`). Directly scanned via `.Scan()`.
- **`db.Exec()`**: Used for statements that do not return rows (`INSERT`, `UPDATE`, `DELETE`). Returns `sql.Result` with `.LastInsertId()` and `.RowsAffected()`.

### 4. Parameterized SQL Queries
- Using `?` placeholders prevents SQL injection vulnerabilities by separating user inputs from SQL logic.

### 5. Error Handling (`if err != nil`)
- Go does not use exceptions. Functions return an `error` value which is explicitly checked. If `err != nil`, an error occurred and must be handled.
