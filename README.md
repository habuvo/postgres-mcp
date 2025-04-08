# PostgreSQL MCP Server

This is a Model Context Protocol (MCP) server that provides an interface to interact with PostgreSQL databases. It allows you to execute queries, perform database operations, and manage transactions through a standardized API.

## Features

- Execute SQL queries with parameterized inputs
- Perform database operations (INSERT, UPDATE, DELETE)
- Execute transactions with multiple statements
- Retrieve database schema information
- Secure parameter handling to prevent SQL injection

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/postgres-mcp.git
   cd postgres-mcp
   ```

2. Install dependencies:

   ```
   go mod download
   ```

3. Configure your database connection:

   ```
   cp .env.example .env
   ```

   Edit the `.env` file with your PostgreSQL database credentials.

4. Build the server:

   ```
   go build -o postgres-mcp
   ```

5. Run the server:

   ```
   ./postgres-mcp
   ```

## Usage

The server exposes the following MCP tools:

### 1. execute_tool

Executes a single SQL statement.

**Request:**

```json
{
  "statement": "INSERT INTO users(name, email) VALUES($1, $2)",
  "arguments": ["Jane Doe", "jane@example.com"]
}
```

**Response:**

```json
{
  "rows affected": 1
}
```

### 2. query_tool

Executes a SELECT query and returns the results.

**Request:**

```json
{
  "statement": "SELECT * FROM users WHERE id = $1",
  "arguments": [1]
}
```

**Response:**

```json
{
  "rows": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  ],
  "columns": ["id", "name", "email"],
  "rowCount": 1
}
```

### 3. transaction_tool

Executes multiple statements in a transaction.

**Request:**

```json
{
  "statements": [
    "INSERT INTO users(name, email) VALUES($1, $2)",
    "INSERT INTO user_roles(user_id, role) VALUES($1, $2)"
  ],
  "arguments": [
    ["Alice", "alice@example.com"],
    [1, "admin"]
  ]
}
```

**Response:**

```json
{
  "results": [
    {
      "statement": 0,
      "rowsAffected": 1
    },
    {
      "statement": 1,
      "rowsAffected": 1
    }
  ]
}
```

### 4. schema_tool

Retrieves schema information for a table.

**Request:**

```json
{
  "table_name": "users"
}
```

**Response:**

```json
{
  "table": "users",
  "columns": [
    {
      "name": "id",
      "type": "integer",
      "nullable": false,
      "default": "nextval('users_id_seq'::regclass)"
    },
    {
      "name": "name",
      "type": "character varying",
      "nullable": false,
      "default": null
    },
    {
      "name": "email",
      "type": "character varying",
      "nullable": false,
      "default": null
    }
  ]
}
```

## Client Example

Here's an example of how to use the MCP client to interact with the server:

```go
package main

import (
 "encoding/json"
 "fmt"
 "log"

 "github.com/mark3labs/mcp-go/client"
)

func main() {
 // Create a new MCP client
 c := client.NewClient("http://localhost:8080")

 // Prepare a query request
 queryReq := map[string]interface{}{
  "statement": "SELECT * FROM users WHERE id = $1",
  "arguments": []interface{}{1},
 }

 // Execute the query
 resp, err := c.Execute("query_tool", queryReq)
 if err != nil {
  log.Fatalf("Failed to execute query: %v", err)
 }

 // Print the response
 prettyJSON, _ := json.MarshalIndent(resp.Data, "", "  ")
 fmt.Println(string(prettyJSON))
}
```

## Security Considerations

- Always use parameterized queries to prevent SQL injection attacks
- Limit database user permissions to only what is necessary
- Use SSL/TLS for database connections in production
- Implement proper authentication and authorization for the MCP server

## License

This project is licensed under the MIT License.
