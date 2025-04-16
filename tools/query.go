package tools

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterQueryTool(s *server.MCPServer, db *sql.DB) {
	queryTool := mcp.NewTool("query_tool",
		mcp.WithDescription("Make query"),
		mcp.WithString("statement",
			mcp.Required(),
			mcp.Description("Statement to be executed"),
		),
		mcp.WithArray("arguments",
			mcp.Required(),
			mcp.Description("Arguments for the statent provided"),
		),
	)

	s.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleQuery(ctx, request, db)
	})
}

func handleQuery(ctx context.Context, req mcp.CallToolRequest, db *sql.DB) (*mcp.CallToolResult, error) {
	statement, ok := req.Params.Arguments["statement"].(string)
	if !ok {
		return nil, fmt.Errorf("statement should be a string")
	}

	args, ok := req.Params.Arguments["arguments"].([]any)
	if !ok {
		return nil, fmt.Errorf("arguments should be an array of arguments")
	}

	// Execute the query with parameters
	rows, err := db.Query(statement, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	// Prepare result
	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the result into the values
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Create a map for this row
		rowMap := make(map[string]interface{})

		for i, col := range columns {
			val := values[i]
			// Handle null values
			if val == nil {
				rowMap[col] = nil
			} else {
				// Try to convert to appropriate type
				switch v := val.(type) {
				case []byte:
					// Convert []byte to string
					rowMap[col] = string(v)
				default:
					rowMap[col] = v
				}
			}
		}

		results = append(results, rowMap)
	}

	// Check for errors after iterating through rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	// Set response data
	return mcp.NewToolResultText(fmt.Sprintf("rows: %s columns: %s, rowCount: %d", results, columns, len(results))), nil
}
