package routenames

const (
	RouteNameLogin                   = "login"
	RouteNameLoginSubmit             = "login.submit"
	RouteNameLogout                  = "logout"
	RouteNameContact                 = "contact"
	RouteNameContactSubmit           = "contact.submit"
	RouteNameAboutUs                 = "about"
	RouteNameHomeInternet            = "home_internet"
	RouteNameEnterpriseInternet      = "enterprise_internet"
	RouteNameLandingPage             = "landing_page"
	RouteNameBlog                    = "blog"
	RouteNameBlogPost                = "blog.post"
	RouteNamePreferences             = "preferences"
	RouteNameGetPhone                = "phone.get"
	RouteNameUpdatePhoneNum          = "phone.save"
	RouteNameGetDisplayName          = "display_name.get"
	RouteNameUpdateDisplayName       = "display_name.save"
	RouteNameGetPhoneVerification    = "phone.verification"
	RouteNameSubmitPhoneVerification = "phone.verification.submit"

	RouteNameNotifications            = "normalNotifications"
	RouteNameNormalNotificationsCount = "normalNotificationsCount"

	RouteNamePrivacyPolicy      = "privacy_policy"
	RouteNameTermsAndConditions = "terms_and_conditions"
	RouteNameRefundPolicy       = "refund_policy"

	RouteNameProfile    = "profile"
	RouteNameInstallApp = "install_app"

	RouteNameMarkNotificationsAsRead    = "markNormalNotificationRead"
	RouteNameMarkAllNotificationsAsRead = "normalNotificationsMarkAllAsRead"

	RouteNameRealtime = "realtime"

	RouteNameFinishOnboarding = "finish_onboarding"
	RouteNameGetBio           = "profileBio.get"
	RouteNameUpdateBio        = "profileBio.post"
	RouteNameProfileBio       = "profileBio"

	RouteNameGetPushSubscriptions             = "push_subscriptions.get"
	RouteNameRegisterSubscription             = "notification_subscriptions.register"
	RouteNameDeleteSubscription               = "notification_subscriptions.delete"
	RouteNameDeleteEmailSubscriptionWithToken = "email_subscriptions.delete_with_token"

	RouteNamePaymentProcessorGetPublicKey = "payment_processor.get_public_key"
	RouteNameCreateCheckoutSession        = "stripe.create_checkout_session"
	RouteNameCreatePortalSession          = "stripe.create_portal_session"
	RouteNamePaymentProcessorWebhook      = "stripe.webhook"
	RouteNamePricingPage                  = "pricing_page"
	RouteNamePaymentProcessorSuccess      = "stripe.success"

	// ISP Features
	RouteNameTicketCreate = "ticket.create"
	RouteNameTicketSubmit = "ticket.submit"
	RouteNameRenewPackage = "package.renew"
	RouteNameChangePlan   = "package.change"
	RouteNameAddFunds     = "balance.add"
)
