package tools

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterSchemaTool(s *server.MCPServer, dbs map[string]*sql.DB) {
	schemaTool := mcp.NewTool("schema_tool",
		mcp.WithDescription("Inspect table schema on specified database"),
		mcp.WithString("database",
			mcp.Required(),
			mcp.Description("Name of the database to inspect"),
		),
		mcp.WithString("table_name",
			mcp.Required(),
			mcp.Description("The name of the table to inspect"),
		))

	s.AddTool(schemaTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleSchema(ctx, request, dbs)
	})
}

// Execute handles the schema request
func handleSchema(ctx context.Context, request mcp.CallToolRequest, dbs map[string]*sql.DB) (*mcp.CallToolResult, error) {
	database, ok := request.Params.Arguments["database"].(string)
	if !ok {
		return nil, fmt.Errorf("database parameter should be a string")
	}

	db, exists := dbs[database]
	if !exists {
		return nil, fmt.Errorf("database '%s' not found in available connections", database)
	}

	tableName, ok := request.Params.Arguments["table_name"].(string)
	if !ok || tableName == "" {
		return nil, fmt.Errorf("table_name must be a non-empty string")
	}

	// Query to get table schema
	query := `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute schema query: %v", err)
	}
	defer rows.Close()

	// Prepare result
	var columns []map[string]any

	for rows.Next() {
		var (
			columnName, dataType, isNullable string
			columnDefault                    sql.NullString
		)

		if err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		column := map[string]interface{}{
			"name":     columnName,
			"type":     dataType,
			"nullable": isNullable == "YES",
		}

		if columnDefault.Valid {
			column["default"] = columnDefault.String
		} else {
			column["default"] = nil
		}

		columns = append(columns, column)
	}

	// Check for errors after iterating through rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	// Set response data
	return mcp.NewToolResultText(fmt.Sprintf("table: %s columns: %s", tableName, columns)), nil
}
