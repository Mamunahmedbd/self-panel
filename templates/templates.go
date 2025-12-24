package templates

type (
	Page string
)

const (
	PageAbout                  Page = "about"
	PageLanding                Page = "landing"
	PageContact                Page = "contact"
	PageBlog                   Page = "blog"
	PageBlogPost               Page = "blog_post"
	PageError                  Page = "error"
	PageNotFound               Page = "not_found"
	PageLogin                  Page = "login"
	PageEmailSubscribe         Page = "email-subscribe"
	PageProfile                Page = "profile"
	PagePhoneNumber            Page = "profile.phone"
	PageDisplayName            Page = "profile.display_name"
	PageInstallApp             Page = "install_app"
	PageDashboard              Page = "dashboard"
	PageNotifications          Page = "notifications"
	PageHealthcheck            Page = "healthcheck"
	PagePricing                Page = "pricing"
	PageSuccessfullySubscribed Page = "successfully_subscribed"
	PagePrivacyPolicy          Page = "privacy_policy"
	PageTermsAndConditions     Page = "terms_and_conditions"
	PageRefundPolicy           Page = "refund_policy"
	PageWiki                   Page = "wiki"

	SSEAnsweredByFriend Page = "sse_answered_by_friend"
)
