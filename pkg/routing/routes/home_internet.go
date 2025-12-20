package routes

import (
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"

	"github.com/labstack/echo/v4"
)

type homeInternet struct {
	ctr controller.Controller
}

func NewHomeInternetRoute(ctr controller.Controller) homeInternet {
	return homeInternet{
		ctr: ctr,
	}
}

func (c *homeInternet) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.LandingPage
	page.Name = templates.PageLanding // Or create a new constant
	page.Title = "Home Internet - High Speed Fiber"
	page.Component = pages.HomeInternet(&page)

	return c.ctr.RenderPage(ctx, page)
}
