package tools

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterExecuteTool(s *server.MCPServer, dbs map[string]*sql.DB) {
	executeTool := mcp.NewTool("execute_tool",
		mcp.WithDescription("Execute SQL statement on specified database"),
		mcp.WithString("database",
			mcp.Required(),
			mcp.Description("Name of the database to execute on"),
		),
		mcp.WithString("statement",
			mcp.Required(),
			mcp.Description("SQL statement(s) to execute. For multiple statements, separate them with semicolons (;)"),
		),
		mcp.WithArray("arguments",
			mcp.Required(),
			mcp.Description("Arguments for the statement provided"),
		),
	)

	s.AddTool(executeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleExecute(ctx, request, dbs)
	})
}

func handleExecute(ctx context.Context, req mcp.CallToolRequest, dbs map[string]*sql.DB) (*mcp.CallToolResult, error) {
	database, ok := req.Params.Arguments["database"].(string)
	if !ok {
		return nil, fmt.Errorf("database parameter should be a string")
	}

	db, exists := dbs[database]
	if !exists {
		return nil, fmt.Errorf("database '%s' not found in available connections", database)
	}

	statement, ok := req.Params.Arguments["statement"].(string)
	if !ok {
		return nil, fmt.Errorf("statement should be a string")
	}

	args, ok := req.Params.Arguments["arguments"].([]any)
	if !ok {
		return nil, fmt.Errorf("arguments should be an array of arguments")
	}

	// Execute the statement
	result, err := db.ExecContext(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	// Get rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("rows affected: %d", rowsAffected)), nil
}
