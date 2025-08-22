"""
MCP Server - Python Implementation
"""

import os
import json
import requests
from pathlib import Path
from typing import Annotated
from pydantic import Field
from mcp.server.fastmcp import FastMCP

# Create MCP server instance
mcp = FastMCP("MCP Server")

def get_config():
    """Get configuration from environment or config file."""
    class Config:
        def __init__(self):
            self.base_url = os.getenv("API_BASE_URL")
            self.bearer_token = os.getenv("API_BEARER_TOKEN")
            
            # Try to load from config file if env vars not set
            if not self.base_url or not self.bearer_token:
                config_path = Path.home() / ".api" / "config.json"
                if config_path.exists():
                    with open(config_path, 'r') as f:
                        config_data = json.load(f)
                        self.base_url = self.base_url or config_data.get("baseURL")
                        self.bearer_token = self.bearer_token or config_data.get("bearerToken")
    
    return Config()

# Add configuration resource
@mcp.resource("config://settings")
def get_config_resource() -> str:
    """Get current configuration settings."""
    config = get_config()
    return json.dumps({
        "base_url": config.base_url,
        "bearer_token": "***" if config.bearer_token else None
    }, indent=2)

# Tool functions
@mcp.tool()
def get_search_text(query: Annotated[str, Field(description="The query! The query is the main parameter used to define a query. You can type several terms separated by spaces or phrases wrapped inside quote ‘”’ characters. For every term, you can also use ‘+’ and ‘-‘ modifier characters to indicate that a term is “mandatory” or “prohibited” (by default, terms are considered to be “mandatory”). For example, in a query such as query=term_a -term_b, sounds including term_b will not match the search criteria. The query does a weighted search over some sound properties including sound tags, the sound name, its description, pack name and the sound id. Therefore, searching for query=123 will find you sounds with id 1234, sounds that have 1234 in the description, in the tags, etc. You’ll find some examples below. Using an empty query (query= or query=\"\\")], filter: Annotated[str, Field(description="Allows filtering query results. See below for more information.")], sort: Annotated[str, Field(description="Indicates how query results should be sorted. See below for a list of the sorting options. By default `sort=score`. <p> <table> <tr> <th>Option</th> <th>Explanation</th> </tr> <tr> <td>score</td> <td>Sort by a relevance score returned by our search engine (default).</td> </tr> <tr> <td>duration_desc <td>Sort by the duration of the sounds, longest sounds first. </tr> <tr> <td>duration_asc <td>Same as above, but shortest sounds first. </tr> <tr> <td>created_desc <td>Sort by the date of when the sound was added. newest sounds first. </tr> <tr> <td>created_asc <td>Same as above, but oldest sounds first. </tr> <tr> <td>downloads_desc <td>Sort by the number of downloads, most downloaded sounds first. </tr> <tr> <td>downloads_asc <td>Same as above, but least downloaded sounds first. </tr> <tr> <td>rating_desc <td>Sort by the average rating given to the sounds, highest rated first. </tr> <tr> <td>rating_asc <td>Same as above, but lowest rated sounds first. </tr> </table> </p>")], group_by_pack: Annotated[str, Field(description="This parameter represents a boolean option to indicate whether to collapse results belonging to sounds of the same pack into single entries in the results list. If `group_by_pack=1` and search results contain more than one sound that belongs to the same pack, only one sound for each distinct pack is returned (sounds with no packs are returned as well). However, the returned sound will feature two extra properties to access these other sounds omitted from the results list: `n_from_same_pack`: indicates how many other results belong to the same pack (and have not been returned) `more_from_same_pack`: uri pointing to the list of omitted sound results of the same pack (also including the result which has already been returned). See examples below. By default `group_by_pack=0`.")], page: Annotated[str, Field(description="Query results are paginated, this parameter indicates what page should be returned. By default `page=1`.")], page_size: Annotated[str, Field(description="Indicates the number of sounds per page to include in the result. By default `page_size=15`, and the maximum is `page_size=150`. Not that with bigger `page_size`, more data will need to be transferred.")]) -> str:
    """Search sounds"""
    try:
        config = get_config()
        
        if not config.base_url or not config.bearer_token:
            return "Error: Missing API configuration. Please set API_BASE_URL and API_BEARER_TOKEN environment variables."
        
        # Build request parameters
        params = {}
        if query: params["query"] = query
        if filter: params["filter"] = filter
        if sort: params["sort"] = sort
        if group_by_pack: params["group_by_pack"] = group_by_pack
        if page: params["page"] = page
        if page_size: params["page_size"] = page_size
        
        # Make API call
        url = f"{config.base_url}/api/unknown"
        
        headers = {
            "Authorization": f"Bearer {config.bearer_token}",
            "Accept": "application/json",
            "Content-Type": "application/json"
        }
        
        response = requests.get(url, headers=headers, params=params)
        
        # Handle HTTP errors
        if response.status_code >= 400:
            try:
                error_data = response.json()
                return f"Failed to format JSON: {json.dumps(error_data, indent=2)}"
            except json.JSONDecodeError:
                return f"Failed to format JSON: {response.text}"
        
        # Parse response
        try:
            result = response.json()
            return json.dumps(result, indent=2)
        except json.JSONDecodeError:
            # Fallback to raw text if JSON parsing fails
            return response.text
            
    except requests.exceptions.ConnectionError as e:
        return f"Request failed: Connection error - {str(e)}"
    except requests.exceptions.Timeout as e:
        return f"Request failed: Request timeout - {str(e)}"
    except requests.exceptions.RequestException as e:
        return f"Request failed: {str(e)}"
    except Exception as e:
        return f"Unexpected error: {str(e)}"


@mcp.tool()
def get_sounds_sound_id(soundId: Annotated[str, Field(description="ID of the sound that needs to be fetched")]) -> str:
    """Details of a sound"""
    try:
        config = get_config()
        
        if not config.base_url or not config.bearer_token:
            return "Error: Missing API configuration. Please set API_BASE_URL and API_BEARER_TOKEN environment variables."
        
        # Build request parameters
        params = {}
        if soundId: params["soundId"] = soundId
        
        # Make API call
        url = f"{config.base_url}/api/unknown"
        
        headers = {
            "Authorization": f"Bearer {config.bearer_token}",
            "Accept": "application/json",
            "Content-Type": "application/json"
        }
        
        response = requests.get(url, headers=headers, params=params)
        
        # Handle HTTP errors
        if response.status_code >= 400:
            try:
                error_data = response.json()
                return f"Failed to format JSON: {json.dumps(error_data, indent=2)}"
            except json.JSONDecodeError:
                return f"Failed to format JSON: {response.text}"
        
        # Parse response
        try:
            result = response.json()
            return json.dumps(result, indent=2)
        except json.JSONDecodeError:
            # Fallback to raw text if JSON parsing fails
            return response.text
            
    except requests.exceptions.ConnectionError as e:
        return f"Request failed: Connection error - {str(e)}"
    except requests.exceptions.Timeout as e:
        return f"Request failed: Request timeout - {str(e)}"
    except requests.exceptions.RequestException as e:
        return f"Request failed: {str(e)}"
    except Exception as e:
        return f"Unexpected error: {str(e)}"


if __name__ == "__main__":
    mcp.run()
