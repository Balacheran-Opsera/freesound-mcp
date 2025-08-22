package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/freesound/mcp-server/config"
	"github.com/freesound/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func SearchtextHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["query"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("query=%v", val))
		}
		if val, ok := args["filter"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("filter=%v", val))
		}
		if val, ok := args["sort"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("sort=%v", val))
		}
		if val, ok := args["group_by_pack"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("group_by_pack=%v", val))
		}
		if val, ok := args["page"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("page=%v", val))
		}
		if val, ok := args["page_size"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("page_size=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/search/text%s", cfg.BaseURL, queryString)
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
		var result []Sound
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

func CreateSearchtextTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_search_text",
		mcp.WithDescription("Search sounds"),
		mcp.WithString("query", mcp.Description("The query! The query is the main parameter used to define a query. You can type several terms separated by spaces or phrases wrapped inside quote ‘”’ characters. For every term, you can also use ‘+’ and ‘-‘ modifier characters to indicate that a term is “mandatory” or “prohibited” (by default, terms are considered to be “mandatory”). For example, in a query such as query=term_a -term_b, sounds including term_b will not match the search criteria. The query does a weighted search over some sound properties including sound tags, the sound name, its description, pack name and the sound id. Therefore, searching for query=123 will find you sounds with id 1234, sounds that have 1234 in the description, in the tags, etc. You’ll find some examples below. Using an empty query (query= or query=\"\") will return all Freeosund sounds.")),
		mcp.WithString("filter", mcp.Description("Allows filtering query results. See below for more information.")),
		mcp.WithString("sort", mcp.Description("Indicates how query results should be sorted. See below for a list of the sorting options. By default `sort=score`. <p> <table>\n  <tr>\n    <th>Option</th>\n    <th>Explanation</th>\n  </tr>\n  <tr>\n    <td>score</td>\n    <td>Sort by a relevance score returned by our search engine (default).</td>\n  </tr>\n  <tr>\n    <td>duration_desc\n    <td>Sort by the duration of the sounds, longest sounds first.\n  </tr>\n  <tr>\n    <td>duration_asc\n    <td>Same as above, but shortest sounds first.\n  </tr>\n  <tr>\n    <td>created_desc\n    <td>Sort by the date of when the sound was added. newest sounds first.\n  </tr>\n  <tr>\n    <td>created_asc\n    <td>Same as above, but oldest sounds first.\n  </tr>\n  <tr>\n    <td>downloads_desc\n    <td>Sort by the number of downloads, most downloaded sounds first.\n  </tr>\n  <tr>\n    <td>downloads_asc\n    <td>Same as above, but least downloaded sounds first.\n  </tr>\n  <tr>\n    <td>rating_desc\n    <td>Sort by the average rating given to the sounds, highest rated first.\n  </tr>\n  <tr>\n    <td>rating_asc\n    <td>Same as above, but lowest rated sounds first.\n  </tr>\n</table> </p>")),
		mcp.WithNumber("group_by_pack", mcp.Description("This parameter represents a boolean option to indicate whether to collapse results belonging to sounds of the same pack into single entries in the results list. If `group_by_pack=1` and search results contain more than one sound that belongs to the same pack, only one sound for each distinct pack is returned (sounds with no packs are returned as well). However, the returned sound will feature two extra properties to access these other sounds omitted from the results list: `n_from_same_pack`: indicates how many other results belong to the same pack (and have not been returned) `more_from_same_pack`: uri pointing to the list of omitted sound results of the same pack (also including the result which has already been returned). See examples below. By default `group_by_pack=0`.")),
		mcp.WithNumber("page", mcp.Description("Query results are paginated, this parameter indicates what page should be returned. By default `page=1`.")),
		mcp.WithNumber("page_size", mcp.Description("Indicates the number of sounds per page to include in the result. By default `page_size=15`, and the maximum is `page_size=150`. Not that with bigger `page_size`, more data will need to be transferred.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    SearchtextHandler(cfg),
	}
}
