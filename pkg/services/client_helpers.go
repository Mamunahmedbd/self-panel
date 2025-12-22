package services

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/pkg/context"
)

// GetAuthenticatedClient retrieves the authenticated client from the session
// Returns nil if no client is authenticated or if there's an error
func (c *Container) GetAuthenticatedClient(ctx echo.Context) (*ent.ClientUser, error) {
	sess, err := session.Get("session", ctx)
	if err != nil {
		return nil, err
	}

	clientID, ok := sess.Values[context.ClientIDKey].(int)
	if !ok {
		return nil, nil // No client ID in session
	}

	// Load the client from database
	client, err := c.ORM.ClientUser.Get(ctx.Request().Context(), clientID)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetClientID retrieves just the client ID from the session
// Returns 0 if no client is authenticated
func (c *Container) GetClientID(ctx echo.Context) int {
	sess, err := session.Get("session", ctx)
	if err != nil {
		return 0
	}

	clientID, ok := sess.Values[context.ClientIDKey].(int)
	if !ok {
		return 0
	}

	return clientID
}

// GetClientUsername retrieves the client username from the session
// Returns empty string if no client is authenticated
func (c *Container) GetClientUsername(ctx echo.Context) string {
	sess, err := session.Get("session", ctx)
	if err != nil {
		return ""
	}

	username, ok := sess.Values[context.ClientUsernameKey].(string)
	if !ok {
		return ""
	}

	return username
}

// IsClientAuthenticated checks if a client is currently authenticated
func (c *Container) IsClientAuthenticated(ctx echo.Context) bool {
	return c.GetClientID(ctx) != 0
}
