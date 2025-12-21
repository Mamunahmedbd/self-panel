package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"
)

type (
	terms struct {
		ctr controller.Controller
	}
)

func NewTermsRoute(ctr controller.Controller) terms {
	return terms{
		ctr: ctr,
	}
}

func (c *terms) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.LandingPage
	page.Name = templates.PageTermsAndConditions
	page.Title = "Terms & Conditions"
	page.HTMX.Request.Boosted = true
	page.Component = pages.Terms(&page)

	return c.ctr.RenderPage(ctx, page)
}
