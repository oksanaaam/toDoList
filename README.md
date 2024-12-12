# ToDo List API

This is a simple API for managing todo lists, implemented in Go, using PostgreSQL as a database.

## Description

The API allows you to create, get, update and delete tasks in a todo list. We use **Gin** as a framework for creating the API and **pgx** for working with PostgreSQL.

### Main functionalities:
- **GET /todos** – get all tasks.
- **GET /todos/:id** – get a task by ID.
- **POST /todos** – add a new task.
- **PUT /todos/:id** – update a task.
- **DELETE /todos/:id** – delete a task.

## Technologies
- **Go** (Golang)
- **Gin** (web framework)
- **PostgreSQL** (database)
- **pgx** (library for working with PostgreSQL)

## Installation and launch

### 1. Cloning the repository
First you need to clone the repository:

```bash
git clone https://github.com/your-username/todolist-api.git
cd todolist-api
```

### 2. Setting up PostgreSQL

1. **Create the database**:
   Use pgAdmin or the command line to create a database:
   ```sql
   CREATE DATABASE todolist;

2. **Create the todos table: Run the following SQL query to create the table**:
    ```sql
    CREATE TABLE todos (
        id VARCHAR PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        status VARCHAR(50) NOT NULL
    );

3. **Set up environment variables**:
    Create a .env file or set the following environment variables in your system:
        - DB_CONNECTION_STRING=postgres://username:password@localhost:5432/todolist
        - Replace username and password with your actual PostgreSQL credentials.

4. **Installing dependencies**:
    Run the following command to install all required Go dependencies:
    `go mod tidy`

5. **Running the project**:
    To start the project, use the following command:
    `go run main.go`

    This will start the server, and the API will be available at http://localhost:8080.

6. **Testing the endpoints**:
    You can test the endpoints using tools like Postman or curl.

    Example requests:

    **GET /todos**:
    GET http://localhost:8080/todos

    **GET /todos/:id**:
    GET http://localhost:8080/todos/1

    **POST /todos**:
    POST http://localhost:8080/todos
    Content-Type: application/json
    Body: {
    "id": "3",
    "title": "New Task",
    "status": "pending"
    }

    **PUT /todos/:id**:
    PUT http://localhost:8080/todos/3
    Content-Type: application/json
    Body: {
    "title": "Updated Task",
    "status": "completed"
    }
    
    **DELETE /todos/:id**:
    DELETE http://localhost:8080/todos/3

