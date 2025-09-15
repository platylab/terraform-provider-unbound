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
			// Get ID from the key
			id, err := strconv.Atoi(key)
			if err != nil {
				log.Printf("Failed to convert ID '%s' to integer: %v", key, err)
				continue
			}

			// Get value
			valueStr, ok := value.(string)
			if !ok {
				log.Printf("Value for ID '%s' is not a string", key)
				continue
			}

			// Split the value into Name and Type
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


// CreateLocalZone - Returns specific local zone
func (c *Client) CreateLocalZone(ctx context.Context, body PostConfigClauseAttributeValueIdJSONRequestBody) (LocalZone, error) {
	// Parameters
	clause    := "server"
	attribute := "local-zone"
	valueId   := "*"

	// Perform a POST request
	resp, err := c.PostConfigClauseAttributeValueId(ctx, clause, attribute, valueId, body)
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
	if !ok || int(status) != 201 {
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
			// Get ID from the key
			id, err := strconv.Atoi(key)
			if err != nil {
				log.Printf("Failed to convert ID '%s' to integer: %v", key, err)
				continue
			}

			// Get value
			valueStr, ok := value.(string)
			if !ok {
				log.Printf("Value for ID '%s' is not a string", key)
				continue
			}

			// Split the value into Name and Type
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


// UpdateLocalZone - Returns specific local zone
func (c *Client) UpdateLocalZone(ctx context.Context, valueId int, body PutConfigClauseAttributeValueIdJSONRequestBody) (LocalZone, error) {
	// Parameters
	clause    := "server"
	attribute := "local-zone"

	// Perform a PUT request
	resp, err := c.PutConfigClauseAttributeValueId(ctx, clause, attribute, strconv.Itoa(valueId), body)
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
		for key, item := range items {
			// Access individual item and assert it's a map
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				return localZone, fmt.Errorf("item '%s' is not a valid map", key)
			}

			// Extract the ID field
			id, err := strconv.Atoi(itemMap["id"].(string))
			if err != nil {
				return localZone, fmt.Errorf("failed to convert ID '%s' to integer: %v", key, err)
			}

			// Extract the "new_value" field
			valueStr, ok := itemMap["new_value"].(string)
			if !ok {
				return localZone, fmt.Errorf("'new_value' field for item '%s' is not a valid string", key)
			}

			// Split the value into Name and Type
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

			// Since we only need the first item, return immediately
			return localZone, nil
		}
	}

	// If no valid items were found (loop completed without returning)
	return LocalZone{}, fmt.Errorf("Failed to parse any valid items in the response")

}

// DeleteLocalZone - Returns specific local zone
func (c *Client) DeleteLocalZone(ctx context.Context, valueId int) (LocalZone, error) {
	// Parameters
	clause    := "server"
	attribute := "local-zone"

	// Perform a DELETE request
	resp, err := c.DeleteConfigClauseAttributeValueId(ctx, clause, attribute, strconv.Itoa(valueId))
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
		for key, item := range items {
			// Access individual item and assert it's a map
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				return localZone, fmt.Errorf("item '%s' is not a valid map", key)
			}

			// Extract the ID field
			id, err := strconv.Atoi(itemMap["id"].(string))
			if err != nil {
				return localZone, fmt.Errorf("failed to convert ID '%s' to integer: %v", key, err)
			}

			// Extract the "old_value" field
			valueStr, ok := itemMap["old_value"].(string)
			if !ok {
				return localZone, fmt.Errorf("'old_value' field for item '%s' is not a valid string", key)
			}

			// Split the value into Name and Type
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

			// Since we only need the first item, return immediately
			return localZone, nil
		}
	}

	// If no valid items were found (loop completed without returning)
	return LocalZone{}, fmt.Errorf("Failed to parse any valid items in the response")

}
