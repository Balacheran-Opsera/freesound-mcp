package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/freesound/mcp-server/config"
	"github.com/freesound/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetsoundbyidHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		soundIdVal, ok := args["soundId"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: soundId"), nil
		}
		soundId, ok := soundIdVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: soundId"), nil
		}
		url := fmt.Sprintf("%s/sounds/%s", cfg.BaseURL, soundId)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// No authentication required for this endpoint
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.Sound
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGetsoundbyidTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_sounds_soundId",
		mcp.WithDescription("Details of a sound"),
		mcp.WithNumber("soundId", mcp.Required(), mcp.Description("ID of the sound that needs to be fetched")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetsoundbyidHandler(cfg),
	}
}
