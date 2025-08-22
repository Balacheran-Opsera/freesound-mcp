package main

import (
	"github.com/freesound/mcp-server/config"
	"github.com/freesound/mcp-server/models"
	tools_search "github.com/freesound/mcp-server/tools/search"
	tools_sound "github.com/freesound/mcp-server/tools/sound"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_search.CreateSearchtextTool(cfg),
		tools_sound.CreateGetsoundbyidTool(cfg),
	}
}
