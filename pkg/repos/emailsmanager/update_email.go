package emailsmanager

import (
	"context"
	"errors"
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/notification"
	"github.com/mikestefanello/pagoda/ent/notificationpermission"
	"github.com/mikestefanello/pagoda/ent/profile"
	"github.com/mikestefanello/pagoda/ent/sentemail"
	"github.com/mikestefanello/pagoda/pkg/domain"
	"github.com/mikestefanello/pagoda/pkg/routing/routenames"
	"github.com/mikestefanello/pagoda/pkg/services"
	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates/emails"
	"github.com/mikestefanello/pagoda/templates/layouts"
)

type UpdateEmailSender struct {
	orm *ent.Client
	// TODO: feels kinda weird pass container here, but refactor later.
	container *services.Container
}

func NewUpdateEmailSender(orm *ent.Client, container *services.Container) *UpdateEmailSender {
	return &UpdateEmailSender{
		orm:       orm,
		container: container,
	}
}

// GetAudience returns the audience who should receive update emails.
// How many days in a row a person with "partner notifs" turned on will receive an email is determined
// by the function that consumes GetAudience. The ideais to only send them an email if they have a new notif in the last n days.
func (e *UpdateEmailSender) GetAudience(ctx context.Context) ([]int, error) {

	oneDayAgo := time.Now().Add(-24 * time.Hour)

	// We don't want to send more than 1 email a day to any user.
	alreadySentEmailFilter := profile.Not(
		profile.HasSentEmailsWith(
			sentemail.And(
				sentemail.CreatedAtGTE(oneDayAgo),
				sentemail.Or(
					sentemail.TypeEQ(sentemail.TypePartnerActivity),
				),
			),
		),
	)

	// Get all users who gave permission for partner updates
	profilesWithOnlyPartnerUpdates, err := e.orm.Profile.Query().
		Where(
			profile.HasNotificationPermissionsWith(
				notificationpermission.PermissionEQ(notificationpermission.Permission(domain.NotificationPermissionNewFriendActivity.Value)),
			),
			// Unread notifications of the below type qualify as "partner updates".
			profile.HasNotificationsWith(
				notification.ReadEQ(false),
				notification.TypeIn(
					notification.Type(domain.NotificationTypeConnectionEngagedWithQuestion.Value),
				),
			),
			alreadySentEmailFilter,
		).
		Select(profile.FieldID).
		All(ctx)

	if err != nil {
		return nil, err
	}

	profileIDs := mapset.NewSet[int]()

	for _, profile := range profilesWithOnlyPartnerUpdates {
		profileIDs.Add(profile.ID)
	}

	profileIDsSlice := profileIDs.ToSlice()

	return profileIDsSlice, nil
}

func (e *UpdateEmailSender) PrepareAndSendUpdateEmailForAll(ctx context.Context) error {
	// profileIDs, err := e.GetAudience(ctx)
	// if err != nil {
	// 	return err
	// }

	// entProfiles, err := e.orm.Profile.
	// 	Query().
	// 	Where(
	// 		profile.IDIn(profileIDs...),
	// 	).
	// 	WithNotifications(func(nq *ent.NotificationQuery) {
	// 		nq.Where(
	// 			notification.ReadEQ(false),
	// 			notification.TypeIn(
	// 				notification.Type(domain.NotificationTypeCommittedRelationshipRequest.Value),
	// 				notification.Type(domain.NotificationTypeConnectionReactedToAnswer.Value),
	// 				notification.Type(domain.NotificationTypeConnectionRequestAccepted.Value),
	// 				notification.Type(domain.NotificationTypeConnectionEngagedWithQuestion.Value),
	// 				notification.Type(domain.NotificationTypeMutualQuestionAnswered.Value),
	// 			),
	// 		)
	// 	}).
	// 	WithNotificationPermissions(func(npq *ent.NotificationPermissionQuery) {
	// 		npq.Where(notificationpermission.PlatformEQ(notificationpermission.Platform(domain.NotificationPlatformEmail.Value)))
	// 	}).
	// 	WithUser(func(uq *ent.UserQuery) {
	// 		uq.Select(user.FieldEmail, user.FieldName)
	// 	}).
	// 	All(ctx)

	// for _, entProfile := range entProfiles {
	// 	err = e.PrepareAndSendUpdateEmailForProfile(ctx, entProfile)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (e *UpdateEmailSender) SendUpdateEmail(
	ctx context.Context, profileID int, email, selfName, partnerName string,
	questionsAnswered, questionsNotAnswered []types.QuestionInEmail,
	dailyUpdatePermissionToken, partnerUpdatePermissionToken string, numPartnerNotifications int,
) error {

	if len(questionsAnswered) > 0 && partnerUpdatePermissionToken == "" {
		return errors.New("partner update permission token missing")
	}

	title := "New questions to answer!"
	questionsAnsweredByFriendButNotSelfTitle := ""

	// TODO: if partner in committed mode has answers waiting, use their name.
	if len(questionsAnswered) > 0 {
		answerText := "answers"
		if len(questionsAnswered) == 1 {
			answerText = "answer"
		}

		questionText := "questions"
		if len(questionsAnswered) == 1 {
			questionText = "question"
		}

		if partnerName == "" {
			title = fmt.Sprintf(
				"%d new %s you can read as soon as you've answered them yourself!", len(questionsAnswered), answerText)
			questionsAnsweredByFriendButNotSelfTitle = fmt.Sprintf(
				"You can read %d new %s from your matches for the following:", len(questionsAnswered), answerText)
		} else {
			title = fmt.Sprintf(
				"%s answered %d new %s you can read as soon as you've answered them yourself!", partnerName, len(questionsAnswered), questionText)
			questionsAnsweredByFriendButNotSelfTitle = fmt.Sprintf(
				"Christa answered %d %s you have not answered:", len(questionsAnswered), questionText)

		}
	}

	url := e.container.Web.Reverse(routenames.RouteNameDeleteEmailSubscriptionWithToken,
		domain.NotificationPermissionNewFriendActivity.Value, partnerUpdatePermissionToken)
	unsubscribePartnerActivityLink := fmt.Sprintf("%s%s", e.container.Config.HTTP.Domain, url)

	data := types.EmailUpdate{
		SelfName:                                 selfName,
		AppName:                                  string(e.container.Config.App.Name),
		SupportEmail:                             e.container.Config.Mail.FromAddress,
		Domain:                                   e.container.Config.HTTP.Domain,
		PartnerName:                              partnerName,
		NumNewNotifications:                      numPartnerNotifications,
		QuestionsAnsweredByFriendButNotSelfTitle: questionsAnsweredByFriendButNotSelfTitle,
		NumQuestionsAnsweredByFriendButNotSelf:   len(questionsAnswered),
		QuestionsAnsweredByFriendButNotSelf:      questionsAnswered,
		QuestionsNotAnsweredInSocialCircle:       questionsNotAnswered,
		UnsubscribePartnerActivityLink:           unsubscribePartnerActivityLink,
	}

	err := e.container.Mail.
		Compose().
		To(email).
		Subject(title).
		TemplateLayout(layouts.Email).
		Component(emails.EmailUpdate(data)).
		Send(ctx)
	if err != nil {
		return err
	}

	emailType := domain.NotificationPermissionNewFriendActivity
	_, err = e.orm.SentEmail.Create().
		SetProfileID(profileID).
		SetType(sentemail.Type(emailType.Value)).
		Save(ctx)

	return err
}
