package routes

import (
	"fmt"
	"regexp"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/packageplan"
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

	// Fetch packages from database
	pkgPlans, err := c.ctr.Container.ORM.PackagePlan.
		Query().
		Where(
			packageplan.IsActive(true),
			packageplan.PriceGT(0),
		).
		Order(ent.Asc(packageplan.FieldPrice)).
		All(ctx.Request().Context())

	if err != nil {
		return err
	}

	var packages []types.InternetPackage
	speedRe := regexp.MustCompile(`(\d+)\s*Mbps`)

	for _, p := range pkgPlans {
		// Extract speed from name
		var speedInt int
		speed := "High"
		if matches := speedRe.FindStringSubmatch(p.Name); len(matches) > 1 {
			speed = matches[1]
			fmt.Sscanf(speed, "%d", &speedInt)
		}

		// Generate features based on speed ranges (Mbps)
		var features []string

		switch {
		case speedInt < 10:
			features = []string{"Unlimited Data", "Bufferless Browsing", "Standard Support", "Shared IP"}
		case speedInt < 12:
			features = []string{"Unlimited Data", "Bufferless Facebook", "Standard Support", "HD Youtube"}
		case speedInt < 15:
			features = []string{"Unlimited Data", "Bufferless Youtube", "24/7 Support", "Multi-Device"}
		case speedInt < 20:
			features = []string{"Unlimited Data", "Full HD Streaming", "Online Class Ready", "Gaming Optimization"}
		case speedInt < 25:
			features = []string{"Unlimited Data", "4K Youtube", "Priority Support", "Public DNS", "Lag-Free Gaming"}
		case speedInt < 30:
			features = []string{"Unlimited Data", "4K Netflix", "Smarthome Ready", "Low Latency", "Pro Gaming"}
		case speedInt < 50:
			features = []string{"Unlimited Data", "Multiple 4K Streams", "Cloud Sync", "Public IP (Optional)", "4K Everything"}
		case speedInt < 100:
			features = []string{"Unlimited Data", "8K Streaming", "Home Office Pro", "Real IP Included", "VIP Support"}
		default: // 100+
			features = []string{"Unlimited Data", "Gigabit Experience", "Enterprise Grade", "Dedicated Manager", "All Features Unlocked"}
		}

		packages = append(packages, types.InternetPackage{
			Name:     p.Name,
			Speed:    speed,
			Price:    fmt.Sprintf("%.0f", p.Price),
			Popular:  speedInt >= 20 && speedInt < 50, // Popular range logic
			Features: features,
		})
	}

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
		Packages: packages,
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
