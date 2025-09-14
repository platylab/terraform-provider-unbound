package unboundapiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// GetLocalData - Returns specific local data
func (c *Client) GetLocalData(ctx context.Context, valueId int) (LocalData, error) {
	// Parameters
	clause := "server"
	attribute := "local-data"

	// Perform a GET request
	resp, err := c.GetConfigClauseAttributeValueId(ctx, clause, attribute, strconv.Itoa(valueId))
	if err != nil {
		log.Fatalf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respJson, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the JSON response into a map
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respJson, &responseMap); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check the "status" field
	status, ok := responseMap["status"].(float64) // JSON numbers are parsed as float64
	if !ok || int(status) != 200 {
    // Extract "error" and "reason" fields as strings with type assertions
    errMessage, _ := responseMap["error"].(string)   // Use type assertion for error message
    errReason, _ := responseMap["reason"].(string)   // Use type assertion for reason

    // Return a formatted error
    return LocalData{}, fmt.Errorf("%s: %s", errMessage, errReason)
	}

	// Check for the "items" field
	items, ok := responseMap["items"].(map[string]interface{})
	if !ok || len(items) == 0 {
		return LocalData{}, fmt.Errorf("No items found in the response")
	}

	// Directly extract the first "item" into a LocalData struct
	var localData LocalData
	if items, ok := responseMap["items"].(map[string]interface{}); ok {
		for key, value := range items {
			// Convert the key (ID) from string to int
			id, err := strconv.Atoi(key)
			if err != nil {
				log.Printf("Failed to convert ID '%s' to integer: %v", key, err)
				continue
			}

			valueStr, ok := value.(string)
			if !ok {
				log.Printf("Value for ID '%s' is not a string", key)
				continue
			}

			parts := strings.SplitN(valueStr, " ", 4) // Split by spaces in 4 parts
			if len(parts) != 4 {
				log.Printf("Failed to parse Entry for ID '%s'", key)
				continue
			}

			// Populate the LocalData struct
			localData = LocalData{
				Id:      id,
				Domain:  strings.Trim(parts[0], `"`), // Remove surrounding `"`,
				Type:    parts[2],
				Value:   strings.Trim(parts[3], `"`),
			}
			
			return localData, nil
		}
	}

	// If no valid items were found (loop completed without returning)
	return LocalData{}, fmt.Errorf("Failed to parse any valid items in the response")

}

