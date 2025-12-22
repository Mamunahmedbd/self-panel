# Quick Reference: PPPoE Client Login

## 🚀 Quick Start

### 1. Create Test Client in Database

```sql
INSERT INTO clients (
    name, username, password, mobile_number, email,
    status, c_name, vendor_id, created_by, created_date
) VALUES (
    'John Doe', 'johndoe', 'test123', '01712345678', 'john@example.com',
    'active', 'MyISP', 1, 'admin', NOW()
);
```

### 2. Start Application

```bash
go run ./cmd/web/main.go
```

### 3. Test Login

- URL: `http://localhost:8000/user/login`
- Username: `johndoe` (or email: `john@example.com`)
- Password: `test123`

## 📖 Common Code Patterns

### Access Client Data in Handler

```go
func MyHandler(c controller.Controller) func(ctx echo.Context) error {
    return func(ctx echo.Context) error {
        // Get authenticated client
        client, err := c.Container.GetAuthenticatedClient(ctx)
        if err != nil {
            return c.Fail(err, "failed to get client")
        }

        if client == nil {
            return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
        }

        // Use client data
        fmt.Printf("Client: %s, Balance: %.2f\n", client.Username, client.Balance)

        return nil
    }
}
```

### Check if Client is Authenticated

```go
if c.Container.IsClientAuthenticated(ctx) {
    // Client is logged in
    username := c.Container.GetClientUsername(ctx)
    fmt.Println("Logged in as:", username)
}
```

### Get Client ID Only

```go
clientID := c.Container.GetClientID(ctx)
if clientID != 0 {
    // Client is authenticated
}
```

### Display Client Info in Template

```go
// In your handler
page := controller.NewPage(ctx)

client, _ := c.Container.GetAuthenticatedClient(ctx)
if client != nil {
    page.Data = map[string]interface{}{
        "ClientName":    client.Name,
        "ClientBalance": client.Balance,
        "ClientPackage": client.PackagePool,
        "ClientStatus":  client.Status,
    }
}

return c.RenderPage(ctx, page)
```

## 🔐 Security Checks

### Verify Account is Active

```go
client, err := c.Container.GetAuthenticatedClient(ctx)
if err != nil {
    return err
}

if client.Status != clientuser.StatusActive {
    return echo.NewHTTPError(
        http.StatusForbidden,
        "Your account is not active",
    )
}
```

### Check Balance Before Action

```go
client, _ := c.Container.GetAuthenticatedClient(ctx)

if client.Balance < requiredAmount {
    return echo.NewHTTPError(
        http.StatusPaymentRequired,
        "Insufficient balance",
    )
}
```

## 🎯 Route Examples

### Create Client Dashboard Route

```go
// In routes/client_dashboard.go
package routes

import (
    "github.com/labstack/echo/v4"
    "github.com/mikestefanello/pagoda/pkg/controller"
    "github.com/mikestefanello/pagoda/templates/pages"
)

type clientDashboard struct {
    ctr controller.Controller
}

func NewClientDashboardRoute(ctr controller.Controller) clientDashboard {
    return clientDashboard{ctr: ctr}
}

func (c *clientDashboard) Get(ctx echo.Context) error {
    page := controller.NewPage(ctx)
    page.Name = "Client Dashboard"
    page.Title = "My Dashboard"

    // Get client data
    client, err := c.ctr.Container.GetAuthenticatedClient(ctx)
    if err != nil {
        return c.ctr.Fail(err, "failed to get client")
    }

    page.Data = map[string]interface{}{
        "Client": client,
    }

    page.Component = pages.ClientDashboard(&page)
    return c.ctr.RenderPage(ctx, page)
}
```

### Register Route

```go
// In routes/router.go, inside coreAuthRoutes function
clientDashboard := NewClientDashboardRoute(ctr)
onboardedGroup.GET("/dashboard", clientDashboard.Get).Name = "client_dashboard"
```

## 📊 Database Queries

### Get Client with Related Data

```go
client, err := c.Container.ORM.ClientUser.
    Query().
    Where(clientuser.IDEQ(clientID)).
    Only(ctx.Request().Context())
```

### Update Client Balance

```go
client, err := c.Container.ORM.ClientUser.
    UpdateOneID(clientID).
    SetBalance(newBalance).
    Save(ctx.Request().Context())
```

### Get All Active Clients

```go
clients, err := c.Container.ORM.ClientUser.
    Query().
    Where(clientuser.StatusEQ(clientuser.StatusActive)).
    All(ctx.Request().Context())
```

### Search Clients

```go
clients, err := c.Container.ORM.ClientUser.
    Query().
    Where(
        clientuser.Or(
            clientuser.NameContains(searchTerm),
            clientuser.UsernameContains(searchTerm),
            clientuser.EmailContains(searchTerm),
        ),
    ).
    Limit(20).
    All(ctx.Request().Context())
```

## 🎨 Template Examples

### Display Client Info

```templ
// In templates/pages/client_dashboard.templ
package pages

import "github.com/mikestefanello/pagoda/pkg/controller"

templ ClientDashboard(page *controller.Page) {
    <div class="container mx-auto p-6">
        if client, ok := page.Data["Client"].(*ent.ClientUser); ok {
            <h1>Welcome, { client.Name }!</h1>
            <div class="grid grid-cols-3 gap-4 mt-6">
                <div class="card">
                    <h3>Balance</h3>
                    <p class="text-2xl">৳{ fmt.Sprintf("%.2f", client.Balance) }</p>
                </div>
                <div class="card">
                    <h3>Package</h3>
                    <p>{ client.PackagePool }</p>
                </div>
                <div class="card">
                    <h3>Status</h3>
                    <p class={ client.Status == "active" ? "text-green-500" : "text-red-500" }>
                        { string(client.Status) }
                    </p>
                </div>
            </div>
        }
    </div>
}
```

## 🔄 Session Management

### Store Additional Data

```go
sess, _ := session.Get("session", ctx)
sess.Values["last_package_view"] = packageID
sess.Values["preferred_language"] = "en"
sess.Save(ctx.Request(), ctx.Response())
```

### Retrieve Session Data

```go
sess, _ := session.Get("session", ctx)
if packageID, ok := sess.Values["last_package_view"].(int); ok {
    // Use packageID
}
```

## 🐛 Debugging

### Log Client Activity

```go
client, _ := c.Container.GetAuthenticatedClient(ctx)
if client != nil {
    ctx.Logger().Infof(
        "Client %s (%d) accessed %s",
        client.Username,
        client.ID,
        ctx.Request().URL.Path,
    )
}
```

### Check Session Contents

```go
sess, _ := session.Get("session", ctx)
ctx.Logger().Debugf("Session values: %+v", sess.Values)
```

## 📝 Validation Examples

### Validate Package Change

```go
func (c *packageChange) Post(ctx echo.Context) error {
    client, _ := c.ctr.Container.GetAuthenticatedClient(ctx)

    // Check if client can change package
    if !client.AutoRenew {
        return echo.NewHTTPError(
            http.StatusBadRequest,
            "Auto-renew must be enabled to change package",
        )
    }

    // Process package change
    // ...
}
```

### Validate Payment

```go
func (c *payment) Post(ctx echo.Context) error {
    client, _ := c.ctr.Container.GetAuthenticatedClient(ctx)

    var form PaymentForm
    ctx.Bind(&form)

    if form.Amount <= 0 {
        return echo.NewHTTPError(
            http.StatusBadRequest,
            "Invalid amount",
        )
    }

    if form.Amount > client.Balance {
        return echo.NewHTTPError(
            http.StatusBadRequest,
            "Insufficient balance",
        )
    }

    // Process payment
    // ...
}
```

## 🎯 Common Tasks

### Task 1: Show Package Details

```go
client, _ := c.Container.GetAuthenticatedClient(ctx)
packageName := client.PackagePool

// Query package details from packages table
// Display to user
```

### Task 2: Top Up Balance

```go
client, _ := c.Container.GetAuthenticatedClient(ctx)
newBalance := client.Balance + topUpAmount

client, err := c.Container.ORM.ClientUser.
    UpdateOneID(client.ID).
    SetBalance(newBalance).
    Save(ctx.Request().Context())
```

### Task 3: Change Package

```go
client, _ := c.Container.GetAuthenticatedClient(ctx)

client, err := c.Container.ORM.ClientUser.
    UpdateOneID(client.ID).
    SetNextUserProfile(newPackageProfile).
    Save(ctx.Request().Context())
```

## 🔗 Useful Links

- Main Docs: `docs/CLIENT_LOGIN_IMPLEMENTATION.md`
- Summary: `IMPLEMENTATION_SUMMARY.md`
- Test: `test_login_logic.go`

## 💡 Tips

1. **Always check for nil**: Client might not be authenticated
2. **Use helper functions**: `GetAuthenticatedClient()` is your friend
3. **Validate status**: Check `client.Status` before sensitive operations
4. **Log activities**: Track client actions for security
5. **Handle errors**: Always check error returns

---

**Quick Reference Version**: 1.0.0
**Last Updated**: 2025-12-22
