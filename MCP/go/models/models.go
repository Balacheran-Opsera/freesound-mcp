package models

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

type Tool struct {
	Definition mcp.Tool
	Handler    func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// Sound represents the Sound schema from the OpenAPI specification
type Sound struct {
	Url string `json:"url,omitempty"` // The URI for this sound on the Freesound website.
	Id int64 `json:"id,omitempty"` // The sound’s unique identifier.
	Name string `json:"name,omitempty"` // The name user gave to the sound.
}
