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
	PageLogin                  Page = "login"
	PageEmailSubscribe         Page = "email-subscribe"
	PagePreferences            Page = "preferences"
	PagePhoneNumber            Page = "preferences.phone"
	PageDisplayName            Page = "preferences.display_name"
	PageInstallApp             Page = "install_app"
	PageProfile                Page = "profile"
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
