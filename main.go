package main

import (
	"log"

	"github.com/habuvo/postgres-mcp/config"
	"github.com/habuvo/postgres-mcp/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load configurations
	configs, err := config.LoadConfigs("POSTGRES_DBS")
	if err != nil {
		log.Fatalf("Failed to load configurations: %v", err)
	}

	// Connect to all databases
	dbs, err := config.ConnectDBs(configs)
	if err != nil {
		log.Printf("Warning: some database connections failed: %v", err)
		if len(dbs) == 0 {
			log.Fatal("No database connections established")
		}
	}

	// Close all connections on exit
	defer func() {
		for _, db := range dbs {
			db.Close()
		}
	}()

	log.Printf("Successfully connected to %d PostgreSQL database(s)!", len(dbs))

	// Initialize the MCP server
	s := server.NewMCPServer(
		"Postgre DB MCP Server",
		"1.0.0",
		server.WithLogging(),
	)

	// Register the PostgreSQL tools with all databases
	tools.RegisterExecuteTool(s, dbs)
	tools.RegisterQueryTool(s, dbs)
	tools.RegisterSchemaTool(s, dbs)
	tools.RegisterTransactionTool(s, dbs)

	// Start the server
	log.Println("Starting MCP server...")

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
