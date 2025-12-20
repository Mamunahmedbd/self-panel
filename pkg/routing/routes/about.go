package routes

import (
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"

	"github.com/labstack/echo/v4"
)

type (
	about struct {
		ctr controller.Controller
	}
)

func NewAboutUsRoute(ctr controller.Controller) about {
	return about{
		ctr: ctr,
	}
}

func (c *about) Get(ctx echo.Context) error {

	page := controller.NewPage(ctx)
	page.Layout = layouts.LandingPage
	page.Name = templates.PageAbout
	page.Title = "About Our Network"
	page.Data = types.AboutData{
		SupportEmail: c.ctr.Container.Config.App.SupportEmail,
		Title:        "Connecting Bangladesh to the Future",
		Subtitle:     "Your Trusted Gateway to the Digital World",
		Description:  "Since 2010, we've been on a mission to provide ultra-fast, reliable, and affordable internet services. Our fiber-optic network spans across the nation, empowering homes and businesses with world-class connectivity.",
		Stats: []types.AboutStat{
			{Value: "15+", Label: "Years Experience"},
			{Value: "50k+", Label: "Happy Clients"},
			{Value: "24/7", Label: "Expert Support"},
			{Value: "99.9%", Label: "Network Uptime"},
		},
		Mission: "To bridge the digital divide by delivering high-speed, secure, and accessible internet solutions that empower every home and business in the country.",
		Vision:  "To be the nationwide leader in digital connectivity, known for our technological excellence and unwavering commitment to customer satisfaction.",
		Features: []types.AboutFeature{
			{
				Title:       "Triple-Redundant Fiber",
				Description: "Our backbone is built on multiple undersea cables to ensure you never lose connection.",
				Icon:        `<svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>`,
			},
			{
				Title:       "Low-Latency Routing",
				Description: "Direct peering with global content providers (Google, Akamai, Facebook) for lag-free experience.",
				Icon:        `<svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>`,
			},
			{
				Title:       "Local Tech Support",
				Description: "Our engineers are located right in your city, ready to visit your site whenever needed.",
				Icon:        `<svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z"/></svg>`,
			},
		},
	}
	page.Component = pages.About(&page)
	page.HTMX.Request.Boosted = true

	return c.ctr.RenderPage(ctx, page)
}
