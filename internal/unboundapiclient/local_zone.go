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

// GetLocalZone - Returns specific local zone
func (c *Client) GetLocalZone(ctx context.Context, valueId int) (LocalZone, error) {
	// Parameters
	clause := "server"
	attribute := "local-zone"

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
    return LocalZone{}, fmt.Errorf("%s: %s", errMessage, errReason)
	}

	// Check for the "items" field
	items, ok := responseMap["items"].(map[string]interface{})
	if !ok || len(items) == 0 {
		return LocalZone{}, fmt.Errorf("No items found in the response")
	}

	// Directly extract the first "item" into a LocalZone struct
	var localZone LocalZone
	if items, ok := responseMap["items"].(map[string]interface{}); ok {
		for key, value := range items {
			// Convert the key (ID) from string to int
			id, err := strconv.Atoi(key)
			if err != nil {
				log.Printf("Failed to convert ID '%s' to integer: %v", key, err)
				continue
			}

			// Split the value into Name and Type
			valueStr, ok := value.(string)
			if !ok {
				log.Printf("Value for ID '%s' is not a string", key)
				continue
			}

			parts := strings.SplitN(valueStr, " ", 2) // Split by first space
			if len(parts) != 2 {
				log.Printf("Failed to parse Name and Type for ID '%s'", key)
				continue
			}

			// Populate the LocalZone struct
			localZone = LocalZone{
				Id:   id,
				Name: strings.Trim(parts[0], `"`), // Remove surrounding `"`,
				Type: parts[1],
			}
			
			return localZone, nil
		}
	}

	// If no valid items were found (loop completed without returning)
	return LocalZone{}, fmt.Errorf("Failed to parse any valid items in the response")

}

