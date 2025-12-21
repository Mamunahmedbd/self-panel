package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gorilla/sessions"
	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/user"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	// authSessionName stores the name of the session which contains authentication data
	authSessionName = "ua"

	// authSessionKeyUserID stores the key used to store the user ID in the session
	authSessionKeyUserID = "user_id"

	// authSessionKeyAuthenticated stores the key used to store the authentication status in the session
	authSessionKeyAuthenticated = "authenticated"
)

// NotAuthenticatedError is an error returned when a user is not authenticated
type NotAuthenticatedError struct{}

// Error implements the error interface.
func (e NotAuthenticatedError) Error() string {
	return "user not authenticated"
}


// AuthClient is the client that handles authentication requests
type AuthClient struct {
	config *config.Config
	orm    *ent.Client
}

// NewAuthClient creates a new authentication client
func NewAuthClient(cfg *config.Config, orm *ent.Client) *AuthClient {
	return &AuthClient{
		config: cfg,
		orm:    orm,
	}
}

// Login logs in a user of a given ID
func (c *AuthClient) Login(ctx echo.Context, userID int) error {

	sess, err := session.Get(authSessionName, ctx)
	if err != nil {
		return err
	}

	sess.Values[authSessionKeyUserID] = userID
	sess.Values[authSessionKeyAuthenticated] = true
	return sess.Save(ctx.Request(), ctx.Response())
}

// Logout logs the requesting user out
func (c *AuthClient) Logout(ctx echo.Context) error {
	sess, err := session.Get(authSessionName, ctx)
	if err != nil {
		return err
	}

	// Overwrite session values
	sess.Values[authSessionKeyAuthenticated] = false

	// TODO: not quite sure why, but resetting the cookie is not needed in the vanilla
	// starter kit from Pagoda. Not sure which one of my changes broke that.
	// Set the cookie to expire immediately
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1, // Set MaxAge to -1 to delete the session
		HttpOnly: true,
	}

	return sess.Save(ctx.Request(), ctx.Response())
}

// GetAuthenticatedUserID returns the authenticated user's ID, if the user is logged in
func (c *AuthClient) GetAuthenticatedUserID(ctx echo.Context) (int, error) {
	sess, err := session.Get(authSessionName, ctx)
	if err != nil {
		return 0, err
	}

	if sess.Values[authSessionKeyAuthenticated] == true {
		return sess.Values[authSessionKeyUserID].(int), nil
	}

	return 0, NotAuthenticatedError{}
}

// GetAuthenticatedUser returns the authenticated user if the user is logged in
func (c *AuthClient) GetAuthenticatedUser(ctx echo.Context) (*ent.User, error) {
	if userID, err := c.GetAuthenticatedUserID(ctx); err == nil {
		return c.orm.User.Query().
			Where(user.ID(userID)).
			WithProfile(func(q *ent.ProfileQuery) {
				q.WithProfileImage(func(pi *ent.ImageQuery) {
					pi.WithSizes(func(s *ent.ImageSizeQuery) {
						s.WithFile()
					})
				})
			}).
			Only(ctx.Request().Context())
	}

	return nil, NotAuthenticatedError{}
}

// SetLastOnlineTimestamp sets the last online time for a user
func (c *AuthClient) SetLastOnlineTimestamp(ctx echo.Context, userID int) error {

	_, err := c.orm.LastSeenOnline.
		Create().
		SetUserID(userID).
		SetSeenAt(time.Now()).
		Save(ctx.Request().Context())

	return err
}

// HashPassword returns a hash of a given password
func (c *AuthClient) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword check if a given password matches a given hash
func (c *AuthClient) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}



// RandomToken generates a random token string of a given length
func (c *AuthClient) RandomToken(length int) (string, error) {
	b := make([]byte, (length/2)+1)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	return token[:length], nil
}

