package tools

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTransactionTool(s *server.MCPServer, dbs map[string]*sql.DB) {
	transactionTool := mcp.NewTool("transaction_tool",
		mcp.WithDescription("Execute queries in a transaction on specified database"),
		mcp.WithString("database",
			mcp.Required(),
			mcp.Description("Name of the database for the transaction"),
		),
		mcp.WithArray("statements",
			mcp.Required(),
			mcp.Description("Statements to be executed in the transaction"),
		),
		mcp.WithArray("arguments",
			mcp.Required(),
			mcp.Description("Arguments for the statements provided"),
		),
	)

	s.AddTool(transactionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleTransaction(ctx, request, dbs)
	})
}

// Execute handles the transaction request.
func handleTransaction(ctx context.Context, request mcp.CallToolRequest, dbs map[string]*sql.DB) (*mcp.CallToolResult, error) {
	database, ok := request.Params.Arguments["database"].(string)
	if !ok {
		return nil, fmt.Errorf("database parameter should be a string")
	}

	db, exists := dbs[database]
	if !exists {
		return nil, fmt.Errorf("database '%s' not found in available connections", database)
	}

	statements, ok := request.Params.Arguments["statements"].([]string)
	if !ok {
		return nil, fmt.Errorf("statements should be an array of strings")
	}

	if len(statements) == 0 {
		return nil, fmt.Errorf("at least one statement is required")
	}

	args, ok := request.Params.Arguments["arguments"].([]any)
	if !ok {
		return nil, fmt.Errorf("arguments should be an array of arguments")
	}

	if len(args) != len(statements) {
		return nil, fmt.Errorf("arguments array should be the same length as statements")
	}

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Prepare response
	results := make([]map[string]any, 0, len(statements))

	// Execute each statement in the transaction
	for i, stmt := range statements {
		// Execute the statement
		result, err := tx.Exec(stmt, args[i])
		if err != nil {
			// Rollback the transaction on error
			tx.Rollback()

			return nil, fmt.Errorf("failed to execute statement %d: %v", i, err)
		}

		// Get rows affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()

			return nil, fmt.Errorf("failed to get rows affected for statement %d: %v", i, err)
		}

		// Add result to results
		results = append(results, map[string]interface{}{
			"statement":    i,
			"rowsAffected": rowsAffected,
		})
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Set response data
	return mcp.NewToolResultText(fmt.Sprintf("%v", results)), nil
}
