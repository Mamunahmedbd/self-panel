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
	blog struct {
		ctr controller.Controller
	}
)

func NewBlogRoute(ctr controller.Controller) blog {
	return blog{
		ctr: ctr,
	}
}

func (c *blog) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.LandingPage
	page.Name = templates.PageBlog
	page.Title = "ISP Insights & Tech Blogs"

	posts := []types.BlogPost{
		{
			ID:       1,
			Title:    "The Future of Fiber Optics: What to Expect in 2026",
			Summary:  "Explore the upcoming breakthroughs in fiber technology that will triple your home internet speeds.",
			Category: "Technology",
			Author:   "Tech Team",
			Date:     "Dec 20, 2025",
			ReadTime: "5 min read",
			ImageURL: "https://images.unsplash.com/photo-1550751827-4bd374c3f58b?q=80&w=800&auto=format&fit=crop",
			Featured: true,
		},
		{
			ID:       2,
			Title:    "How to Optimize Your Home Wi-Fi for Seamless Gaming",
			Summary:  "Tired of lag? Follow our expert guide to positioning your router and managing network traffic.",
			Category: "Guides",
			Author:   "Support Team",
			Date:     "Dec 18, 2025",
			ReadTime: "4 min read",
			ImageURL: "https://images.unsplash.com/photo-1542751371-adc38448a05e?q=80&w=800&auto=format&fit=crop",
		},
		{
			ID:       3,
			Title:    "Why High Latency is Killing Your Business Productivity",
			Summary:  "Learn why ping matters as much as bandwidth for enterprise-level operations.",
			Category: "Enterprise",
			Author:   "Sales Team",
			Date:     "Dec 15, 2025",
			ReadTime: "6 min read",
			ImageURL: "https://images.unsplash.com/photo-1460925895917-afdab827c52f?q=80&w=800&auto=format&fit=crop",
		},
		{
			ID:       4,
			Title:    "Protecting Your Home Network from Modern Security Threats",
			Summary:  "Essential tips to keep your IoT devices and personal data safe from sophisticated cyber attacks.",
			Category: "Security",
			Author:   "Security Ops",
			Date:     "Dec 12, 2025",
			ReadTime: "7 min read",
			ImageURL: "https://images.unsplash.com/photo-1550751827-4bd374c3f58b?q=80&w=800&auto=format&fit=crop",
		},
	}

	page.Data = types.BlogPage{
		Title:      "ISP Insights",
		Subtitle:   "Latest news, tech guides, and networking trends from the front lines of connectivity.",
		Posts:      posts,
		Categories: []string{"All Post", "Technology", "Guides", "Enterprise", "Security"},
	}

	page.Component = pages.Blog(&page)

	return c.ctr.RenderPage(ctx, page)
}
