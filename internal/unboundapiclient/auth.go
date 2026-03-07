package unboundapiclient

import (
	"context"
	"fmt"
	"net/http"
)

// AuthenticatedClient wraps the generated client with authentication functionality.
type AuthenticatedClient struct {
	*ClientWithResponses
	username string
}

// NewAuthenticatedClient creates a new authenticated client.
func NewAuthenticatedClient(server string, username, password string) (*AuthenticatedClient, error) {
	// Create the base client
	client, err := NewClientWithResponses(server)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Authenticate and get token
	resp, err := client.PostLoginWithResponse(context.Background(), PostLoginJSONRequestBody{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Check for successful login
	if resp.JSON200 == nil {
		if resp.JSON400 != nil {
			return nil, fmt.Errorf("login failed: %s: %s", *resp.JSON400.Error, *resp.JSON400.Reason)
		}
		if resp.JSON401 != nil {
			return nil, fmt.Errorf("login failed: %s: %s", *resp.JSON401.Error, *resp.JSON401.Reason)
		}
		return nil, fmt.Errorf("login failed: unexpected status %d", resp.StatusCode())
	}

	// Create request editor to add Authorization header
	token := *resp.JSON200.AccessToken
	c, ok := client.ClientInterface.(*Client)
	if !ok {
		return nil, fmt.Errorf("unexpected client type: %T", client.ClientInterface)
	}
	c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	})

	return &AuthenticatedClient{
		ClientWithResponses: client,
		username:            username,
	}, nil
}

// GetUsername returns the authenticated username.
func (ac *AuthenticatedClient) GetUsername() string {
	return ac.username
}
