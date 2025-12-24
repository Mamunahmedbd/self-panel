package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/ent"
	profilePkg "github.com/mikestefanello/pagoda/ent/profile"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/domain"
	"github.com/mikestefanello/pagoda/pkg/repos/msg"
	"github.com/mikestefanello/pagoda/pkg/routing/routenames"
	routeNames "github.com/mikestefanello/pagoda/pkg/routing/routenames"

	"github.com/mikestefanello/pagoda/pkg/repos/notifierrepo"
	"github.com/mikestefanello/pagoda/pkg/repos/profilerepo"
	"github.com/mikestefanello/pagoda/pkg/repos/subscriptions"
	"github.com/mikestefanello/pagoda/pkg/services"
	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"
	"github.com/rs/zerolog/log"
)

type (
	profilePrefsRoute struct {
		ctr controller.Controller
		orm *ent.Client
	}

	profileBioFormData struct {
		Bio        string `form:"bio" validate:"required"`
		Submission controller.FormSubmission
	}
)

func NewProfilePrefsRoute(ctr controller.Controller, orm *ent.Client) profilePrefsRoute {
	return profilePrefsRoute{
		ctr: ctr,
		orm: orm,
	}
}

func (p *profilePrefsRoute) GetBio(ctx echo.Context) error {
	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	prof := usr.QueryProfile().Select(profilePkg.FieldBio).FirstX(ctx.Request().Context())

	page := controller.NewPage(ctx)
	page.Layout = layouts.Main
	page.Component = pages.About(&page)
	page.Name = templates.PageProfile

	page.Form = &types.ProfileBioFormData{
		Bio: prof.Bio,
	}

	if form := ctx.Get(context.FormKey); form != nil {
		page.Form = form.(*types.ProfileBioFormData)
	}

	return p.ctr.RenderPage(ctx, page)
}

func (p *profilePrefsRoute) UpdateBio(ctx echo.Context) error {
	// Create a new instance of geolocationPoint to hold the incoming data
	var profileBioData types.ProfileBioFormData
	ctx.Set(context.FormKey, &profileBioData)

	if err := ctx.Bind(&profileBioData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid bio data")
	}

	if err := profileBioData.Submission.Process(ctx, profileBioData); err != nil {
		return p.ctr.Fail(err, "unable to process form submission")
	}
	if profileBioData.Submission.HasErrors() {
		return p.GetBio(ctx)
	}

	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

	_, err := p.orm.Profile.UpdateOneID(profile.ID).SetBio(profileBioData.Bio).Save(ctx.Request().Context())
	if err != nil {
		return err
	}

	return p.GetBio(ctx)
}

type profile struct {
	ctr                            controller.Controller
	profileRepo                    profilerepo.ProfileRepo
	pushNotificationsRepo          *notifierrepo.PwaPushNotificationsRepo
	notificationSendPermissionRepo *notifierrepo.NotificationSendPermissionRepo
	subscriptionsRepo              *subscriptions.SubscriptionsRepo
	smsSenderRepo                  *notifierrepo.SMSSender
}

func NewProfileRoute(
	ctr controller.Controller,
	profileRepo *profilerepo.ProfileRepo,
	pushNotificationsRepo *notifierrepo.PwaPushNotificationsRepo,
	notificationSendPermissionRepo *notifierrepo.NotificationSendPermissionRepo,
	subscriptionsRepo *subscriptions.SubscriptionsRepo,
	smsSenderRepo *notifierrepo.SMSSender,
) profile {
	return profile{
		ctr:                            ctr,
		profileRepo:                    *profileRepo,
		pushNotificationsRepo:          pushNotificationsRepo,
		notificationSendPermissionRepo: notificationSendPermissionRepo,
		subscriptionsRepo:              subscriptionsRepo,
		smsSenderRepo:                  smsSenderRepo,
	}
}

func (g *profile) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.Main
	page.Component = pages.Profile(&page)
	page.Name = templates.PageProfile

	var data *types.ProfileSettingsData
	var err error

	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

    // Populate form with current data
    page.Form = &types.UserProfileUpdateForm{
        Name:         usr.Name,
        Email:        usr.Email,
        PhoneNumber:  profile.PhoneNumberE164,
        Description:  profile.Description,
        AddressLine1: profile.AddressLine1,
        AddressLine2: profile.AddressLine2,
        City:         profile.City,
        District:     profile.District,
        Upazila:      profile.Upazila,
        UnionName:    profile.UnionName,
        Zip:          profile.Zip,
    }

    if form := ctx.Get(context.FormKey); form != nil {
		page.Form = form.(*types.UserProfileUpdateForm)
	}

	data, err = g.getCurrPreferencesData(ctx)

	if err != nil {
		return err
	}

	subscribedEndpoints, err := g.pushNotificationsRepo.GetPushSubscriptionEndpoints(ctx.Request().Context(), profile.ID)
	if err != nil {
		return err
	}

	addPushSubscriptionEndpoint := fmt.Sprintf("%s%s",
		g.ctr.Container.Config.HTTP.Domain, ctx.Echo().Reverse(
			routeNames.RouteNameRegisterSubscription, domain.NotificationPlatformPush.Value)) + "?csrf=" + page.CSRF
	deletePushSubscriptionEndpoint := fmt.Sprintf("%s%s",
		g.ctr.Container.Config.HTTP.Domain, ctx.Echo().Reverse(
			routeNames.RouteNameDeleteSubscription, domain.NotificationPlatformPush.Value)) + "?csrf=" + page.CSRF

	addFCMPushSubscriptionEndpoint := fmt.Sprintf("%s%s",
		g.ctr.Container.Config.HTTP.Domain, ctx.Echo().Reverse(
			routeNames.RouteNameRegisterSubscription, domain.NotificationPlatformFCMPush.Value)) + "?csrf=" + page.CSRF
	deleteFCMPushSubscriptionEndpoint := fmt.Sprintf("%s%s",
		g.ctr.Container.Config.HTTP.Domain, ctx.Echo().Reverse(
			routeNames.RouteNameDeleteSubscription, domain.NotificationPlatformFCMPush.Value)) + "?csrf=" + page.CSRF

	addEmailSubscriptionEndpoint := fmt.Sprintf("%s%s",
		g.ctr.Container.Config.HTTP.Domain, ctx.Echo().Reverse(
			routeNames.RouteNameRegisterSubscription, domain.NotificationPlatformEmail.Value)) + "?csrf=" + page.CSRF
	deleteEmailSubscriptionEndpoint := fmt.Sprintf("%s%s",
		g.ctr.Container.Config.HTTP.Domain, ctx.Echo().Reverse(
			routeNames.RouteNameDeleteSubscription, domain.NotificationPlatformEmail.Value)) + "?csrf=" + page.CSRF

	addSmsSubscriptionEndpoint := fmt.Sprintf("%s%s",
		g.ctr.Container.Config.HTTP.Domain, ctx.Echo().Reverse(
			routeNames.RouteNameRegisterSubscription, domain.NotificationPlatformSMS.Value)) + "?csrf=" + page.CSRF
	deleteSmsSubscriptionEndpoint := fmt.Sprintf("%s%s",
		g.ctr.Container.Config.HTTP.Domain, ctx.Echo().Reverse(
			routeNames.RouteNameDeleteSubscription, domain.NotificationPlatformSMS.Value)) + "?csrf=" + page.CSRF

	permissions, err := g.notificationSendPermissionRepo.GetPermissions(ctx.Request().Context(), profile.ID)
	if err != nil {
		return err
	}

	notificationPermissions := types.UserNotificationPermissionsData{
		VapidPublicKey:                g.ctr.Container.Config.App.VapidPublicKey,
		PermissionPartnerActivity:     permissions[domain.NotificationPermissionNewFriendActivity],
		SubscribedEndpoints:           subscribedEndpoints,
		PhoneSubscriptionEnabled:      profile.PhoneNumberE164 != "" && profile.PhoneVerified,
		NotificationTypeQueryParamKey: domain.PermissionNotificationType,

		AddPushSubscriptionEndpoint:    addPushSubscriptionEndpoint,
		DeletePushSubscriptionEndpoint: deletePushSubscriptionEndpoint,

		AddFCMPushSubscriptionEndpoint:    addFCMPushSubscriptionEndpoint,
		DeleteFCMPushSubscriptionEndpoint: deleteFCMPushSubscriptionEndpoint,

		AddEmailSubscriptionEndpoint:    addEmailSubscriptionEndpoint,
		DeleteEmailSubscriptionEndpoint: deleteEmailSubscriptionEndpoint,

		AddSmsSubscriptionEndpoint:    addSmsSubscriptionEndpoint,
		DeleteSmsSubscriptionEndpoint: deleteSmsSubscriptionEndpoint,
	}

	data.NotificationPermissionsData = notificationPermissions

	page.Data = data
	page.HTMX.Request.Boosted = true

	if page.IsFullyOnboarded {
		page.ShowBottomNavbar = true
		page.SelectedBottomNavbarItem = domain.BottomNavbarItemSettings
	}

	return g.ctr.RenderPage(ctx, page)
}

func (g *profile) SaveProfile(ctx echo.Context) error {
	var form types.UserProfileUpdateForm
	ctx.Set(context.FormKey, &form)

	if err := ctx.Bind(&form); err != nil {
		return g.ctr.Fail(err, "unable to parse form")
	}

	if err := form.Submission.Process(ctx, form); err != nil {
		return g.ctr.Fail(err, "unable to process submission")
	}

	if form.Submission.HasErrors() {
		return g.renderProfileDetails(ctx)
	}

	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

	_, err := g.ctr.Container.ORM.User.UpdateOneID(usr.ID).
		SetName(form.Name).
		SetEmail(form.Email).
		Save(ctx.Request().Context())

	if err != nil {
		return err
	}

	_, err = g.ctr.Container.ORM.Profile.UpdateOneID(profile.ID).
		SetDescription(form.Description).
		SetAddressLine1(form.AddressLine1).
		SetAddressLine2(form.AddressLine2).
		SetCity(form.City).
		SetDistrict(form.District).
		SetUpazila(form.Upazila).
		SetUnionName(form.UnionName).
		SetZip(form.Zip).
		SetPhoneNumberE164(form.PhoneNumber).
		Save(ctx.Request().Context())

	if err != nil {
		return err
	}

	msg.Success(ctx, "Profile updated successfully")
	form.Submission.Message = "Profile updated successfully"

	// Re-set the form in context with the success message
	ctx.Set(context.FormKey, &form)

	return g.renderProfileDetails(ctx)
}

func (g *profile) renderProfileDetails(ctx echo.Context) error {
	page := controller.NewPage(ctx)

	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

	// Get the form from context or create a new one with current data
	if form := ctx.Get(context.FormKey); form != nil {
		page.Form = form.(*types.UserProfileUpdateForm)
	} else {
		page.Form = &types.UserProfileUpdateForm{
			Name:         usr.Name,
			Email:        usr.Email,
			PhoneNumber:  profile.PhoneNumberE164,
			Description:  profile.Description,
			AddressLine1: profile.AddressLine1,
			AddressLine2: profile.AddressLine2,
			City:         profile.City,
			District:     profile.District,
			Upazila:      profile.Upazila,
			UnionName:    profile.UnionName,
			Zip:          profile.Zip,
		}
	}

	// Get profile settings data for the read-only stats
	data, err := g.getCurrPreferencesData(ctx)
	if err != nil {
		return err
	}
	page.Data = data

	// Render just the ProfileDetails component without the layout
	return pages.ProfileDetails(&page).Render(ctx.Request().Context(), ctx.Response().Writer)
}

func (g *profile) getCurrPreferencesData(ctx echo.Context) (*types.ProfileSettingsData, error) {

	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)

	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

	// Make sure to check if birthdate is non-nil
	birthdateStr := profile.Birthdate.UTC().Format("2006-01-02")

	activePlan, subscriptionExpiredOn, isTrial, err := g.subscriptionsRepo.GetCurrentlyActiveProduct(
		ctx.Request().Context(), profile.ID,
	)

	if err != nil {
		return nil, err
	}

	data := &types.ProfileSettingsData{
		Bio:                     profile.Bio,
		PhoneNumberInE164Format: profile.PhoneNumberE164,
		CountryCode:             profile.CountryCode,
		SelfBirthdate:           birthdateStr,
		IsProfileFullyOnboarded: profilerepo.IsProfileFullyOnboarded(profile),
		DefaultBio:              domain.DefaultBio,
		DefaultBirthdate:        domain.DefaultBirthdate.Format("2006-01-02"),

        // New fields
        Name:         usr.Name,
        Username:     usr.Username,
        Email:        usr.Email,
        Status:       usr.Status,
        Description:  profile.Description,
        AddressLine1: profile.AddressLine1,
        AddressLine2: profile.AddressLine2,
        City:         profile.City,
        District:     profile.District,
        Upazila:      profile.Upazila,
        UnionName:    profile.UnionName,
        Zip:          profile.Zip,
        CreatedDate:  usr.CreatedAt.Format("2006-01-02 15:04:05"),

		// if IsPaymentsEnabled is true, none of the subscription stuff matters and the entire app will be free
		IsPaymentsEnabled: g.ctr.Container.Config.App.OperationalConstants.PaymentsEnabled,
		IsTrial:           isTrial,
	}

	if activePlan != nil {
		data.ActiveSubscriptionPlan = *activePlan
	} else {
		data.ActiveSubscriptionPlan = domain.ProductTypeFree
	}

	if subscriptionExpiredOn != nil {
		data.MonthlySybscriptionExpiration = subscriptionExpiredOn
	}
	return data, nil
}

func (p *profile) GetPhoneComponent(ctx echo.Context) error {

	page := controller.NewPage(ctx)
	page.Layout = layouts.Main
	page.Component = pages.EditPhonePage(&page)
	page.Name = templates.PageProfile
	page.HTMX.Request.Boosted = true

	return p.ctr.RenderPage(ctx, page)
}

func (p *profile) GetPhoneVerificationComponent(ctx echo.Context) error {
	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

	page := controller.NewPage(ctx)
	page.Layout = layouts.Main
	page.Name = templates.PageProfile
	page.Form = &types.PhoneNumberVerification{}
	page.Component = pages.PhoneVerificationField(&page)
	page.Data = &types.SmsVerificationCodeInfo{
		ExpirationInMinutes: p.ctr.Container.Config.Phone.ValidationCodeExpirationMinutes,
	}

	if form := ctx.Get(context.FormKey); form != nil {
		page.Form = form.(*types.PhoneNumberVerification)
	}

	_, err := p.smsSenderRepo.CreateConfirmationCode(ctx.Request().Context(), profile.ID, profile.PhoneNumberE164)
	if err != nil {
		log.Error().Err(err).Msg("failed to send verification code.")
		msg.Danger(ctx, "Failed to send verification code 😨")
		return p.ctr.RenderPage(ctx, page)
	}

	return p.ctr.RenderPage(ctx, page)
}

func (p *profile) SubmitPhoneVerificationCode(ctx echo.Context) error {

	var form types.PhoneNumberVerification
	ctx.Set(context.FormKey, &form)

	// Parse the form values
	if err := ctx.Bind(&form); err != nil {
		return p.ctr.Fail(err, "unable to parse verification code form")
	}

	if err := form.Submission.Process(ctx, form); err != nil {
		return p.ctr.Fail(err, "unable to process form submission")
	}

	if form.Submission.HasErrors() {
		return p.GetPhoneVerificationComponent(ctx)
	}

	if form.VerificationCode == "" {
		form.Submission.SetFieldError("VerificationCode", "Invalid code")
		msg.Danger(ctx, "Invalid code. Please try again.")
		return p.GetPhoneVerificationComponent(ctx)
	}

	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

	valid, err := p.smsSenderRepo.VerifyConfirmationCode(ctx.Request().Context(), profile.ID, form.VerificationCode)
	if err != nil || !valid {

		form.Submission.SetFieldError("VerificationCode", "Invalid code")
		msg.Danger(ctx, "Invalid code. Please try again.")
		return p.GetPhoneVerificationComponent(ctx)
	}

	page := controller.NewPage(ctx)
	page.Layout = layouts.Main
	page.Name = templates.PageProfile
	page.Form = &types.PhoneNumberVerification{}
	page.Component = pages.PhoneVerificationField(&page)

	msg.Success(ctx, "Success! Your phone number was confirmed.")

	return p.GetPhoneVerificationComponent(ctx)
}

type phoneNumberFormData struct {
	PhoneNumberE164Format string `form:"phone_number_e164" validate:"required"`
	CountryCode           string `form:"country_code" validate:"required"`
	Submission            controller.FormSubmission
}

func (p *profile) SavePhoneInfo(ctx echo.Context) error {
	// Create a new instance of geolocationPoint to hold the incoming data
	var phoneNumberFormData phoneNumberFormData
	ctx.Set(context.FormKey, &phoneNumberFormData)

	if err := ctx.Bind(&phoneNumberFormData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid bio data")
	}

	if err := phoneNumberFormData.Submission.Process(ctx, phoneNumberFormData); err != nil {
		return p.ctr.Fail(err, "unable to process form submission")
	}

	if phoneNumberFormData.Submission.HasErrors() {
		return p.ctr.Redirect(ctx, "profile")
	}

	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

	_, err := p.ctr.Container.ORM.Profile.
		UpdateOneID(profile.ID).
		SetCountryCode(phoneNumberFormData.CountryCode).
		SetPhoneNumberE164(phoneNumberFormData.PhoneNumberE164Format).
		Save(ctx.Request().Context())

	return err
}





type onboarding struct {
	ctr        controller.Controller
	orm        *ent.Client
	taskRunner *services.TaskClient
}

func NewOnboardingRoute(
	ctr controller.Controller, orm *ent.Client, taskRunner *services.TaskClient,
) onboarding {
	return onboarding{
		ctr:        ctr,
		orm:        orm,
		taskRunner: taskRunner,
	}
}

func (p *onboarding) Get(ctx echo.Context) error {
	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	profile := usr.QueryProfile().FirstX(ctx.Request().Context())

	_, err := p.orm.Profile.
		UpdateOneID(profile.ID).
		SetFullyOnboarded(true).
		Save(ctx.Request().Context())
	if err != nil {
		return err
	}

	return p.ctr.RedirectWithDetails(ctx, routenames.RouteNameProfile, "?just_finished_onboarding=true", http.StatusFound)
}
