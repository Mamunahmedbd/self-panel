package routes

import (
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/types"
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

	data := types.ProductLandingPage{
		Hero: types.HeroSection{
			Badge:       "High-Speed Fiber Connectivity",
			Title:       "Home Internet",
			TitleItalic: "Powering Your Digital Life",
			Text:        "Experience the next generation of fiber-to-the-home technology. Ultra-stable, ultra-fast, and built for modern families.",
			PrimaryCTA: types.Link{
				Text: "View Packages",
				URL:  "#packages",
			},
			SecondaryCTA: types.Link{
				Text: "Check Coverage",
				URL:  "/contact",
			},
		},
		Features: []types.Feature{
			{
				Title:       "Ultra Fast",
				Description: "Dedicated fiber bandwidth ensuring consistent speeds even during peak hours.",
			},
			{
				Title:       "Secure Connection",
				Description: "Built-in DDoS protection and enterprise-grade security for your home network.",
			},
			{
				Title:       "24/7 Support",
				Description: "Expert technical support team ready to assist you anytime, day or night.",
			},
		},
		Packages: []types.InternetPackage{
			{
				Name:     "Entry",
				Speed:    "15",
				Price:    "500",
				Popular:  false,
				Features: []string{"Unlimited Data", "2k Youtube Bufferless", "Standard Support"},
			},
			{
				Name:     "Standard",
				Speed:    "30",
				Price:    "800",
				Popular:  true,
				Features: []string{"Unlimited Data", "4k Youtube Bufferless", "Priority Support", "Public DNS"},
			},
			{
				Name:     "Premium",
				Speed:    "50",
				Price:    "1200",
				Popular:  false,
				Features: []string{"Unlimited Data", "4k Everything", "VIP Support", "BDIX Ultra Fast"},
			},
			{
				Name:     "Elite",
				Speed:    "100",
				Price:    "2000",
				Popular:  false,
				Features: []string{"Unlimited Data", "8k Ready", "Personal Account Manager", "Real IP Optional"},
			},
		},
		Technical: types.TechnicalSection{
			Title:       "Technical Excellence in Every Connection",
			Description: "We don't just provide internet; we provide a robust ecosystem designed for the modern digital era. Our PPPoE backbone is optimized for low latency gaming and stutter-free streaming.",
			Stats: []types.Stat{
				{Value: "99.9%", Label: "Uptime Promise"},
				{Value: "<5ms", Label: "Local Latency"},
			},
			ImageURL: "/files/isp_home_internet.png",
			BulletPoints: []string{
				"Dual-Stack IPv4 & IPv6 Support",
				"Advanced BDIX peering for local content",
				"Managed FTTH fiber optic infrastructure",
			},
		},
		FAQs: []types.QAItem{
			{
				Question: "What is PPPoE?",
				Answer:   "Point-to-Point Protocol over Ethernet is the industry standard for home broadband, providing a secure, authenticated connection between your router and our network.",
			},
			{
				Question: "Is there a contract?",
				Answer:   "We offer flexible month-to-month plans. No long-term commitments required for our standard home packages.",
			},
		},
	}

	page.Data = data
	page.HTMX.Request.Boosted = true
	page.Component = pages.HomeInternet(&page)

	return c.ctr.RenderPage(ctx, page)
}
