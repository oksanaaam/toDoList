# ToDo List API

This is a simple API for managing todo lists, implemented in Go, using **PostgreSQL** and **MongoDB** as a database.

## Description

The API allows you to create, get, update and delete tasks in a todo list. We use **Gin** as a framework for creating the API and **pgx** for working with PostgreSQL.

### Main functionalities:
- **GET /** – Gin API homepage.
- **GET /todos** – get all tasks.
- **GET /todos/:id** – get a task by ID.
- **POST /todos** – add a new task.
- **PUT /todos/:id** – update a task.
- **DELETE /todos/:id** – delete a task.

## Technologies
- **Go** (Golang)
- **Gin** (web framework)
- **PostgreSQL** (database)
- **MongoDB** (database)
- **pgx** (library for working with PostgreSQL)

## Installation and launch

### 1. Cloning the repository
First you need to clone the repository:

```bash
git clone https://github.com/your-username/todolist-api.git
cd todolist-api
```

### 2. Setting up PostgreSQL

**PostgreSQL**
1. **Create the database**:
   Use pgAdmin or the command line to create a database:
   ```sql
   CREATE DATABASE todolist;

2. **Create the todos table: Run the following SQL query to create the table**:
    ```sql
    CREATE TABLE todos (
    id VARCHAR PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('created', 'in progress', 'done')),
    image_path VARCHAR(255),
    );

3. **Set up environment variables**:
    Create a .env file or set the following environment variables in your system:
        *DB_TYPE* – choice of database: mongo or postgres
        *MONGO_URI* – URI for connecting to MongoDB (default: mongodb://localhost:27017)
        *MONGO_DB_NAME* – MongoDB database name
        *MONGO_COLLECTION_NAME* – the name of the collection in MongoDB
        *SERVER_ADDRESS* – address for the API (eg localhost:8080)
        *DB_CONNECTION_STRING* – connection string to PostgreSQL (default: postgres://username:password@localhost:5432/todolist)

        Replace username and password with your actual credentials.

4. **Installing dependencies**:
    Run the following command to install all required Go dependencies:
    `go mod tidy`

5. **Running the project**:

**MongoDB**:
   To start MongoDB via Docker, use the command:
   `docker-compose up --build`
   This will build and start the container with MongoDB.  
   After using run:
   `docker-compose down`
   
   To start the project local with Postgres, use the following command:
    `go run cmd/server/main.go`
    This will start the server, and the API will be available at http://localhost:8080.

**Load Balancer**:
If you want to use the Load Balancer to distribute traffic between multiple API servers, make sure to follow these steps:

1. Start Multiple API Servers:  
    Run the following commands in separate terminals for each API server:  
    Server 1 (on localhost:8080):
        `go run cmd/server/main.go`

    Server 2 (on localhost:8081):
        `go run cmd/server/main.go`

    Server 3 (on localhost:8082):
        `go run cmd/server/main.go`
2. Run Load Balancer:
    The load balancer will start automatically when you run the API servers, and it will listen on port 8085.  
    If port 8085 is already in use, the load balancer will not start, and a message will be displayed in the terminal.  

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

