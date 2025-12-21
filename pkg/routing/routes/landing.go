package routes

import (
	"fmt"

	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/routing/routenames"
	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"

	"github.com/labstack/echo/v4"
)

type (
	landingPage struct {
		ctr controller.Controller
	}
)

func NewLandingPageRoute(ctr controller.Controller) landingPage {
	return landingPage{
		ctr: ctr,
	}
}

func (c *landingPage) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.LandingPage

	if page.AuthUser != nil {
		return c.ctr.Redirect(ctx, routenames.RouteNameHomeFeed)

	}

	var data types.LandingPage

	page.Metatags.Description = "Ship your app in record time."
	page.Metatags.Keywords = []string{"Boilerplate", "HTMX", "AlpineJS", "Javascript", "Starter kit", "Startup", "Solopreneur", "Indie Hacking"}
	data = types.LandingPage{
		AppName:           string(c.ctr.Container.Config.App.Name),
		UserSignupEnabled: c.ctr.Container.Config.App.OperationalConstants.UserSignupEnabled,

		Title:      "Experience Ultra-Fast Internet",
		Subtitle:   "Reliable, high-speed connectivity for your home and business. Stream, game, and work without limits.",
		GetNowText: "View Packages",
		IntroTitle: "Why Choose Us?",

		HowItWorksTitle: "Our Promise",

		Quote1: "Since switching to this ISP, my gaming lag is gone and downloads are instant. Best service in the city!",
		Quote2: "Customer support is actually helpful and 24/7. Highly recommended for anyone working from home.",

		ExampleQuestion1: "Fastest Speeds",
		ExampleQuestion2: "Unlimited Data",
		ExampleQuestion3: "24/7 Support",

		AboutUsTitle1: "Connecting you to what matters.",
		AboutUsText1:  "We provide reliable fiber-optic internet services with a focus on speed, uptime, and transparency. No hidden fees, just great internet.",
		AboutUsTitle2: "Our Mission",
		AboutUsText2:  "To bridge the digital divide and provide affordable, high-quality internet access to every household.",

		QAItems: []types.QAItem{
			{
				Question: "How fast is the installation?",
				Answer:   "Typically, we can get you connected within 24-48 hours of your request.",
			},
			{
				Question: "Are there any data caps?",
				Answer:   "No! All our packages come with truly unlimited data.",
			},
			{
				Question: "Can I upgrade my plan later?",
				Answer:   "Yes, you can upgrade your speed instantly through our self-care panel.",
			},
			{
				Question: "What payment methods do you accept?",
				Answer:   "We accept bKash, Credit Cards, and Bank Transfers.",
			},
		},

		BackgroundPhoto2lg: "/files/isp_hero_banner.png", // We will move the generated image here
		BackgroundPhoto2xl: "/files/isp_hero_banner.png",
	}

	data.UserSignupEnabled = c.ctr.Container.Config.App.OperationalConstants.UserSignupEnabledOnLandingPage
	data.ContactEmail = c.ctr.Container.Config.Mail.FromAddress
	data.ProductProCode = c.ctr.Container.Config.App.OperationalConstants.ProductProCode
	data.ProductProPrice = fmt.Sprintf("%.2f", c.ctr.Container.Config.App.OperationalConstants.ProductProPrice)
	data.IsPaymentEnabled = true // Enable visual payment section
	page.Data = data
	page.Name = templates.PageLanding
	page.Component = pages.LandingPage(&page)
	page.HTMX.Request.Boosted = true

	// if c.ctr.Container.Config.App.Environment == config.EnvProduction {
	// 	page.Cache.Enabled = true
	// } else {
	// 	page.Cache.Enabled = false
	// }

	return c.ctr.RenderPage(ctx, page)
}
