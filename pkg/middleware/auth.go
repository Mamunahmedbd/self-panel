package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/repos/profilerepo"
	"github.com/mikestefanello/pagoda/pkg/repos/subscriptions"
	"github.com/mikestefanello/pagoda/pkg/routing/routenames"
	"github.com/mikestefanello/pagoda/pkg/services"
	"github.com/rs/zerolog/log"
)

// LoadAuthenticatedUser loads the authenticated user, if one, and stores in context
func LoadAuthenticatedUser(
	authClient *services.AuthClient, profileRepo *profilerepo.ProfileRepo, subscriptionsRepo *subscriptions.SubscriptionsRepo,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u, err := authClient.GetAuthenticatedUser(c)
			switch err.(type) {
			case *ent.NotFoundError:
				c.Logger().Warn("auth user not found")
			case services.NotAuthenticatedError:
			case nil:
				c.Set(context.AuthenticatedUserKey, u)
				if u.Edges.Profile != nil {
					c.Set(context.ProfileFullyOnboarded, profilerepo.IsProfileFullyOnboarded(u.Edges.Profile))

					// if subscriptionsRepo != nil {
					// 	activeProduct, _, err := subscriptionsRepo.GetCurrentlyActiveProduct(c.Request().Context(), u.Edges.Profile.ID)
					// 	if err != nil {
					// 		log.Error().Err(err).Int("userID", u.ID).Int("profileID", u.Edges.Profile.ID).Msg("failed to get active product in middleware for user")
					// 	}
					// 	c.Set(context.ActiveProductPlan, activeProduct)

					// }
				}
				if profileRepo != nil {
					// TODO: cache profile photo URL somewhere in the stack
					c.Set(context.AuthenticatedUserProfilePicURL, profileRepo.GetProfilePhotoThumbnailURL(u.ID))
				}
				c.Logger().Infof("auth user loaded in to context: %d", u.ID)
			default:
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("error querying for authenticated user: %v", err),
				)
			}

			return next(c)
		}
	}
}


// RequireAuthentication requires that the user be authenticated in order to proceed
func RequireAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get(context.AuthenticatedUserKey) == nil {
				// Get the session
				sess, err := session.Get("session", c)
				if err != nil {
					log.Error().Err(err).Msg("failed to open session to save redirectAfterLogin URL to it")
				} else {
					// Store the original URL they were trying to access
					currentURL := c.Request().RequestURI
					sess.Values["redirectAfterLogin"] = currentURL
					sess.Save(c.Request(), c.Response())

				}

				// Redirect to login page
				url := c.Echo().Reverse(routenames.RouteNameLogin)
				return c.Redirect(http.StatusSeeOther, url)
				// Note: leaving original code commented out in case there are unforeseen consequences...so I remember this change which may have caused it...
				// return echo.NewHTTPError(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// RequireNoAuthentication requires that the user not be authenticated in order to proceed
func RequireNoAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if u := c.Get(context.AuthenticatedUserKey); u != nil {
				url := c.Echo().Reverse(routenames.RouteNameProfile)
				return c.Redirect(http.StatusSeeOther, url)
			}

			return next(c)
		}
	}
}
