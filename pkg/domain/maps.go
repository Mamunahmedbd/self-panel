package domain

// Initialize the map of NotificationPermissionType to NotificationPermission
var NotificationPermissionMap = map[NotificationPermissionType]NotificationPermission{
	NotificationPermissionNewFriendActivity: {
		Title:      "Partner activity",
		Subtitle:   "Answers you missed, sent at most once a day.",
		Permission: NotificationPermissionNewFriendActivity.Value,
	},
}

var NotificationCenterButtonText = map[NotificationType]string{
	NotificationTypeConnectionEngagedWithQuestion: "Answer",
}

// DeleteOnceReadNotificationTypesMap is a map of notification types th;oiSJDfiujladijrgoizdikrjgat can be deleted once seen.
// Note that the boolean doesn't matter, this is just a lazy way of creating a set in Go.
