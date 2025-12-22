package services

import (
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/clienttxn"
	"github.com/mikestefanello/pagoda/ent/packageplan"
	"github.com/mikestefanello/pagoda/ent/radacct"
	"github.com/mikestefanello/pagoda/ent/ticket"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/types"
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

// ... existing GetClientID, GetClientUsername, IsClientAuthenticated functions ...

// GetISPProfileData gathers all data needed for the ISP client dashboard
func (c *Container) GetISPProfileData(ctx echo.Context) (*types.ISPProfileData, error) {
	client, err := c.GetAuthenticatedClient(ctx)
	if err != nil || client == nil {
		return nil, err
	}

	data := &types.ISPProfileData{
		Client:           client,
		Balance:          client.Balance,
		AvailableBalance: client.Balance, // For now, same as balance
		AutoRenew:        client.AutoRenew,
		ValidUntil:       client.PaymentDate,
	}

	// 1. Get current package
	if client.UserProfile != "" {
		pkg, err := c.ORM.PackagePlan.Query().
			Where(packageplan.ProfileNameEQ(client.UserProfile)).
			First(ctx.Request().Context())
		if err == nil {
			data.CurrentPackage = pkg
		}
	}

	// 2. Determine Package Status
	data.PackageStatus = "Active"
	if client.PaymentDate != nil && time.Now().After(*client.PaymentDate) {
		data.PackageStatus = "Expired"
	}
	if client.Status != "active" {
		data.PackageStatus = "Inactive"
	}

	// 3. Get Usage stats
	data.Usage = c.GetUsageStats(ctx, client.Username)

	// 4. Get Payment history
	data.Payments, _ = c.GetPaymentHistory(ctx, client.Username, 10)

	// 5. Get Session history
	data.Sessions, _ = c.GetSessionHistory(ctx, client.Username, 10)

	// 6. Get Recent tickets
	data.Tickets, _ = c.GetRecentTickets(ctx, client.ID, 5)

	return data, nil
}

func (c *Container) GetUsageStats(ctx echo.Context, username string) types.ISPUsageStats {
	var stats types.ISPUsageStats
	dbCtx := ctx.Request().Context()

	// Today's usage
	today := time.Now().Truncate(24 * time.Hour)
	rows, err := c.ORM.RadAcct.Query().
		Where(
			radacct.UsernameEQ(username),
			radacct.AcctstarttimeGTE(today),
		).
		All(dbCtx)
	if err == nil {
		for _, r := range rows {
			if r.Acctinputoctets != nil {
				stats.Today += uint64(*r.Acctinputoctets)
			}
			if r.Acctoutputoctets != nil {
				stats.Today += uint64(*r.Acctoutputoctets)
			}
		}
	}

	// Total usage
	// Better way would be sum in DB, but for now we aggregate
	// (Pagination would be better for high volume)

	return stats
}

func (c *Container) GetPaymentHistory(ctx echo.Context, username string, limit int) ([]*ent.ClientTxn, error) {
	return c.ORM.ClientTxn.Query().
		Where(clienttxn.ClientUsernameEQ(username)).
		Order(ent.Desc(clienttxn.FieldTransactionDate)).
		Limit(limit).
		All(ctx.Request().Context())
}

func (c *Container) GetSessionHistory(ctx echo.Context, username string, limit int) ([]*ent.RadAcct, error) {
	return c.ORM.RadAcct.Query().
		Where(radacct.UsernameEQ(username)).
		Order(ent.Desc(radacct.FieldAcctstarttime)).
		Limit(limit).
		All(ctx.Request().Context())
}

func (c *Container) GetRecentTickets(ctx echo.Context, clientID int, limit int) ([]*ent.Ticket, error) {
	return c.ORM.Ticket.Query().
		Where(ticket.ClientID(clientID)).
		Order(ent.Desc(ticket.FieldCreatedAt)).
		Limit(limit).
		All(ctx.Request().Context())
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
