package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"
)

type (
	refundPolicy struct {
		ctr controller.Controller
	}
)

func NewRefundPolicyRoute(ctr controller.Controller) refundPolicy {
	return refundPolicy{
		ctr: ctr,
	}
}

func (c *refundPolicy) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.LandingPage
	page.Name = templates.PageRefundPolicy
	page.Title = "Refund & Delivery Policy"
	page.HTMX.Request.Boosted = true
	page.Component = pages.RefundPolicy(&page)
	page.Data = types.AboutData{
		SupportEmail: c.ctr.Container.Config.App.SupportEmail,
	}

	return c.ctr.RenderPage(ctx, page)
}
