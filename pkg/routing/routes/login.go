package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/clientuser"
	"github.com/mikestefanello/pagoda/ent/user"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/repos/msg"
	routeNames "github.com/mikestefanello/pagoda/pkg/routing/routenames"

	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type (
	login struct {
		ctr controller.Controller
	}
)

func NewLoginRoute(ctr controller.Controller) login {
	return login{
		ctr: ctr,
	}
}

func (c *login) Get(ctx echo.Context) error {

	page := controller.NewPage(ctx)
	page.Layout = layouts.Auth
	page.Name = templates.PageLogin
	page.Title = "Log in"
	page.Form = &types.LoginForm{}
	page.Component = pages.Login(&page)
	page.HTMX.Request.Boosted = true

	// TODO: below is a bit of a hack. We're sometimes left with a stale CSRF token
	// in the cookies because the user was not actively logged out before their session
	// expired. As a workaround, invalidate any related cookie before attempting to login.
	c.ctr.Container.Auth.Logout(ctx)

	if form := ctx.Get(context.FormKey); form != nil {
		page.Form = form.(*types.LoginForm)
	}

	return c.ctr.RenderPage(ctx, page)
}

func (c *login) Post(ctx echo.Context) error {
	var form types.LoginForm
	ctx.Set(context.FormKey, &form)

	authFailed := func(message string) error {
		if message == "" {
			message = "Invalid username or password. Please try again."
		}
		msg.Danger(ctx, message)
		return c.Get(ctx)
	}

	// Parse the form values
	if err := ctx.Bind(&form); err != nil {
		return c.ctr.Fail(err, "unable to parse login form")
	}

	if err := form.Submission.Process(ctx, form); err != nil {
		return c.ctr.Fail(err, "unable to process form submission")
	}

	if form.Submission.HasErrors() {
		return c.Get(ctx)
	}

	// Trim whitespace from input
	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)

	// Validate input
	if username == "" {
		form.Submission.SetFieldError("Username", "Username is required")
		return c.Get(ctx)
	}

	if password == "" {
		form.Submission.SetFieldError("Password", "Password is required")
		return c.Get(ctx)
	}

	// Attempt to load the client by username only
	client, err := c.ctr.Container.ORM.ClientUser.
		Query().
		Where(clientuser.UsernameEQ(username)).
		Only(ctx.Request().Context())

	switch {
	case ent.IsNotFound(err):
		ctx.Logger().Debugf("client not found: username=%s", username)
		return authFailed("Invalid username or password")
	case err != nil:
		return c.ctr.Fail(err, "error querying client during login")
	}

	// Check if client account is active
	if client.Status != clientuser.StatusActive {
		ctx.Logger().Debugf("client account is not active: username=%s, status=%s", username, client.Status)
		return authFailed("Your account is not active. Please contact support.")
	}

	// Check password (plain-text comparison for PPPoE users)
	if client.Password != password {
		ctx.Logger().Debugf("password incorrect for username=%s", username)
		return authFailed("Invalid username or password")
	}

	// Create or get user account for session management
	// First, check if a user with this email exists
	usr, err := c.ctr.Container.ORM.User.
		Query().
		Where(user.Email(strings.ToLower(client.Email))).
		Only(ctx.Request().Context())

	if ent.IsNotFound(err) {
		// Create a new user account linked to this client
		usr, err = c.ctr.Container.ORM.User.
			Create().
			SetName(client.Name).
			SetEmail(strings.ToLower(client.Email)).
			SetPassword("client-" + client.Username). // Dummy password, not used for client login
			Save(ctx.Request().Context())
		if err != nil {
			return c.ctr.Fail(err, "unable to create user account for client")
		}

		// Create profile for the user
		_, err = c.ctr.Container.ORM.Profile.
			Create().
			SetUser(usr).
			SetBio("PPPoE Client - " + client.Username).
			SetFullyOnboarded(true). // Mark as onboarded
			Save(ctx.Request().Context())
		if err != nil {
			return c.ctr.Fail(err, "unable to create profile for client")
		}
	} else if err != nil {
		return c.ctr.Fail(err, "error querying user during login")
	}

	// Store client ID in session for future reference
	sess, _ := session.Get("session", ctx)
	sess.Values["client_id"] = client.ID
	sess.Values["client_username"] = client.Username
	sess.Save(ctx.Request(), ctx.Response())

	// Log the user in
	err = c.ctr.Container.Auth.Login(ctx, usr.ID)
	if err != nil {
		return c.ctr.Fail(err, "unable to log in user")
	}

	msg.Success(ctx, fmt.Sprintf("Welcome back, <strong>%s</strong>. You are now logged in.", client.Name))

	redirect, err := redirectAfterLogin(ctx)
	if err != nil {
		return err
	}
	if redirect {
		return nil
	}

	// Redirect to profile/dashboard
	return c.ctr.Redirect(ctx, routeNames.RouteNameProfile)
}


// redirectAfterLogin redirects a now logged-in user to a previously requested page.
func redirectAfterLogin(ctx echo.Context) (bool, error) {
	sess, _ := session.Get("session", ctx)

	// Retrieve the redirect URL if it exists
	redirectURL, ok := sess.Values["redirectAfterLogin"].(string)
	if ok && redirectURL != "" {
		// Clear the redirect URL from session
		delete(sess.Values, "redirectAfterLogin")
		sess.Save(ctx.Request(), ctx.Response())

		// Redirect to the originally requested URL
		return true, ctx.Redirect(http.StatusFound, redirectURL)
	}
	return false, nil // Or redirect to a default route if nothing is in the session
}
