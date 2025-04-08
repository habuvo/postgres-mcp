package tools

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterExecuteTool(s *server.MCPServer, db *sql.DB) {
	executeTool := mcp.NewTool("execute_tool",
		mcp.WithDescription("Execute statement"),
		mcp.WithString("statement",
			mcp.Required(),
			mcp.Description("Statement to be executed"),
		),
		mcp.WithArray("arguments",
			mcp.Required(),
			mcp.Description("Arguments for the statement provided"),
		),
	)

	s.AddTool(executeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleExecute(ctx, request, db)
	})
}

func handleExecute(ctx context.Context, req mcp.CallToolRequest, db *sql.DB) (*mcp.CallToolResult, error) {
	statement, ok := req.Params.Arguments["statement"].(string)
	if !ok {
		return nil, fmt.Errorf("statement should be a string")
	}

	args, ok := req.Params.Arguments["arguments"].([]any)
	if !ok {
		return nil, fmt.Errorf("arguments should be an array of arguments")
	}

	// Execute the statement
	result, err := db.Exec(statement, args)
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
