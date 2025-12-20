package routes

import (
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"

	"github.com/labstack/echo/v4"
)

type enterpriseInternet struct {
	ctr controller.Controller
}

func NewEnterpriseInternetRoute(ctr controller.Controller) enterpriseInternet {
	return enterpriseInternet{
		ctr: ctr,
	}
}

func (c *enterpriseInternet) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.LandingPage
	page.Name = templates.PageLanding
	page.Title = "Enterprise Internet - Dedicated Bandwidth & SLA"

	data := types.ProductLandingPage{
		Hero: types.HeroSection{
			Badge:       "Enterprise-Grade Connectivity",
			Title:       "Dedicated Internet Access",
			TitleItalic: "Powering Your Global Ambitions",
			Text:        "Unleash your business potential with high-performance, dedicated fiber connectivity. Guaranteed speeds, 99.99% SLA, and proactive monitoring.",
			PrimaryCTA: types.Link{
				Text: "Get a Quote",
				URL:  "/contact",
			},
			SecondaryCTA: types.Link{
				Text: "View Solutions",
				URL:  "#packages",
			},
		},
		Features: []types.Feature{
			{
				Title:       "Dedicated Bandwidth",
				Description: "1:1 symmetric speeds with zero contention. Your bandwidth is yours alone, 100% of the time.",
			},
			{
				Title:       "Mission-Critical SLA",
				Description: "Industry-leading 99.99% uptime guarantee with aggressive credit-backed service level agreements.",
			},
			{
				Title:       "Proactive Support",
				Description: "24/7/365 Network Operations Center (NOC) monitoring with dedicated account managers.",
			},
		},
		Packages: []types.InternetPackage{
			{
				Name:     "Business Lite",
				Speed:    "50",
				Price:    "5000",
				Popular:  false,
				Features: []string{"Dedicated Fiber", "Symmetric Upload/Download", "99.9% SLA", "8x5 Gold Support"},
			},
			{
				Name:     "Corporate",
				Speed:    "100",
				Price:    "9500",
				Popular:  true,
				Features: []string{"Full Dedicated Access", "Proactive Monitoring", "99.95% SLA", "24/7 Platinum Support", "1 Static IP included"},
			},
			{
				Name:     "Enterprise",
				Speed:    "200",
				Price:    "18000",
				Popular:  false,
				Features: []string{"Premium Routing", "Multi-homed Network", "99.99% SLA", "Dedicated Account Manager", "Block of 8 IPs"},
			},
			{
				Name:     "Custom Fiber",
				Speed:    "1Gbps+",
				Price:    "Custom",
				Popular:  false,
				Features: []string{"Dark Fiber Options", "Direct Cloud Connect", "BGP Peering", "Disaster Recovery", "Custom SLA"},
			},
		},
		Technical: types.TechnicalSection{
			Title:       "Architected for Performance & Reliability",
			Description: "Our enterprise network is built on a redundant fiber ring architecture, ensuring zero single points of failure. We peer directly with global Tier-1 providers for the lowest possible latency.",
			Stats: []types.Stat{
				{Value: "99.99%", Label: "Monthly Uptime"},
				{Value: "<2ms", Label: "IXP Latency"},
			},
			ImageURL: "/files/isp_enterprise_internet.png",
			BulletPoints: []string{
				"Redundant Fiber Entry Points",
				"Carrier-Grade NAT & IPv6 Ready",
				"Layer-2 VPN & MPLS Options Available",
			},
		},
		FAQs: []types.QAItem{
			{
				Question: "What is the difference between Home and Enterprise internet?",
				Answer:   "Enterprise internet provides dedicated, non-contended bandwidth with a guaranteed Service Level Agreement (SLA), whereas home internet is shared among multiple users in an area.",
			},
			{
				Question: "How long does installation take?",
				Answer:   "Standard installation typically takes 7-14 business days, depending on your building's fiber readiness and local permitting.",
			},
		},
	}

	page.Data = data
	page.Component = pages.EnterpriseInternet(&page)

	return c.ctr.RenderPage(ctx, page)
}
