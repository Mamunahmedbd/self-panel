package types

import (
	"time"

	"github.com/mikestefanello/pagoda/pkg/domain"
)

type (
	ProfileSettingsData struct {
		// Form data
		Bio                     string
		PhoneNumberInE164Format string
		CountryCode             string
		SelfBirthdate           string

        // New fields
        Name         string
        Username     string
        Email        string
        Status       string
        Description  string
        AddressLine1 string
        AddressLine2 string
        City         string
        District     string
        Upazila      string
        UnionName    string
        Zip          string
        CreatedDate  string

		// Validation data
		IsProfileFullyOnboarded bool
		DefaultBio              string
		DefaultBirthdate        string

		IsPaymentsEnabled             bool
		ActiveSubscriptionPlan        domain.ProductType
		IsTrial                       bool
		MonthlySybscriptionExpiration *time.Time

		NotificationPermissionsData UserNotificationPermissionsData
	}

    UserProfileUpdateForm struct {
        Name         string `form:"name" validate:"required"`
        Email        string `form:"email" validate:"required,email"`
        PhoneNumber  string `form:"mobile_number"`
        Description  string `form:"description"`
        AddressLine1 string `form:"address_line1"`
        AddressLine2 string `form:"address_line2"`
        City         string `form:"city"`
        District     string `form:"district"`
        Upazila      string `form:"upazila"`
        UnionName    string `form:"union_name"`
        Zip          string `form:"zip"`
        Submission   FormSubmission
    }

	UserNotificationPermissionsData struct {
		// Permissions                    []domain.NotificationPermission
		PermissionPartnerActivity     domain.NotificationPermission
		VapidPublicKey                string
		SubscribedEndpoints           []string
		PhoneSubscriptionEnabled      bool
		NotificationTypeQueryParamKey string

		AddPushSubscriptionEndpoint    string
		DeletePushSubscriptionEndpoint string

		AddFCMPushSubscriptionEndpoint    string
		DeleteFCMPushSubscriptionEndpoint string

		AddEmailSubscriptionEndpoint    string
		DeleteEmailSubscriptionEndpoint string

		AddSmsSubscriptionEndpoint    string
		DeleteSmsSubscriptionEndpoint string
	}
)