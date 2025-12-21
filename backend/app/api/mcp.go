package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
)

func startMCPServer(services *services.AllServices) error {
	s := server.NewMCPServer(
		"Homebox MCP",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Tool: search_items
	searchItemsTool := mcp.NewTool("search_items",
		mcp.WithDescription("Search for items in Homebox"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query string"),
		),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("UUID of the group to search in"),
		),
	)

	s.AddTool(searchItemsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		groupIDStr, err := request.RequireString("group_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		gid, err := uuid.Parse(groupIDStr)
		if err != nil {
			return mcp.NewToolResultError("Invalid group_id UUID"), nil
		}

		res, err := services.Items.Search(ctx, gid, query, 0, 10)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
		}

		// Format result
		var output string
		for _, item := range res.Items {
			locName := "Unknown"
			if item.Location != nil {
				locName = item.Location.Name
			}
			output += fmt.Sprintf("ID: %s | AssetID: %s | Name: %s | Location: %s\n", item.ID, item.AssetID, item.Name, locName)
		}

		if output == "" {
			output = "No items found."
		}

		return mcp.NewToolResultText(output), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "MCP Server error: %v\n", err)
		return err
	}

	return nil
}
