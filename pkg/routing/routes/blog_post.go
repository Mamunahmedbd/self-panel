package routes

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"
)

type (
	blogPost struct {
		ctr controller.Controller
	}
)

func NewBlogPostRoute(ctr controller.Controller) blogPost {
	return blogPost{
		ctr: ctr,
	}
}

func (c *blogPost) Get(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, _ := strconv.Atoi(idParam)

	page := controller.NewPage(ctx)
	page.Layout = layouts.LandingPage
	page.Name = templates.PageBlogPost

	// Dummy data for the single post
	post := types.BlogPost{
		ID:       id,
		Title:    "The Future of Fiber Optics: What to Expect in 2026",
		Summary:  "Explore the upcoming breakthroughs in fiber technology that will triple your home internet speeds.",
		Category: "Technology",
		Author:   "Mamun Ahmed",
		Date:     "Dec 20, 2025",
		ReadTime: "5 min read",
		ImageURL: "https://images.unsplash.com/photo-1550751827-4bd374c3f58b?q=80&w=1200&auto=format&fit=crop",
		Content: `
			<p class="text-lg leading-relaxed mb-6">
				As we move into 2026, the landscape of digital connectivity is undergoing a seismic shift. Fiber optic technology, once a luxury for major enterprise hubs, has become the backbone of the modern digital home. But what lies ahead?
			</p>
			<h2 class="text-3xl font-black text-slate-900 dark:text-white mb-6">The Dawn of Terabit Connectivity</h2>
			<p class="mb-6">
				Upcoming breakthroughs in hollow-core fiber optics are promising to reduce latency by nearly 30% while increasing bandwidth capacity exponentially. Unlike traditional silica-core fibers, light travels through air in hollow-core designs, allowing for near-light-speed data transmission.
			</p>
			<blockquote class="border-l-4 border-indigo-600 pl-6 my-8 italic text-xl text-slate-700 dark:text-gray-300">
				"The next decade of internet speeds won't just be about faster downloads; it will be about the elimination of distance as a barrier to real-time collaboration."
			</blockquote>
			<h2 class="text-3xl font-black text-slate-900 dark:text-white mb-6">Green Networking</h2>
			<p class="mb-6">
				Sustainability is at the core of the 2026 roadmap. Modern Optical Line Terminals (OLTs) are being designed with AI-driven power management that scales energy consumption based on active network demand, reducing the carbon footprint of global ISPs by up to 25%.
			</p>
			<img src="https://images.unsplash.com/photo-1544197150-b99a580bb7a8?q=80&w=1200&auto=format&fit=crop" class="rounded-3xl my-10 shadow-2xl" alt="Network Infrastructure" />
			<h2 class="text-3xl font-black text-slate-900 dark:text-white mb-6">Conclusion</h2>
			<p>
				We are standing at the threshold of a new era. For businesses and families alike, these advancements mean more than just faster Netflix streams—they represent the infrastructure for the metaverse, remote robotic surgery, and truly global, seamless communication.
			</p>
		`,
	}

	page.Title = post.Title
	page.Data = post
	page.Component = pages.BlogPost(&page)

	return c.ctr.RenderPage(ctx, page)
}
