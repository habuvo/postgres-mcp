package main

import (
	"log"

	"github.com/habuvo/postgres-mcp/config"
	"github.com/habuvo/postgres-mcp/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to the PostgreSQL database!")

	// Initialize the MCP server
	s := server.NewMCPServer(
		"Postgre DB MCP Server",
		"1.0.0",
		server.WithLogging(),
	)

	// Register the PostgreSQL tools
	tools.RegisterExecuteTool(s, db)
	tools.RegisterQueryTool(s, db)
	tools.RegisterSchemaTool(s, db)
	tools.RegisterTransactionTool(s, db)

	// Start the server
	log.Println("Starting MCP server...")

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
