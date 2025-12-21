package tasks

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/notification"
)

// ----------------------------------------------------------

// ----------------------------------------------------------
const TypeDeleteStaleNotifications = "notification.recycling"

type (
	DeleteStaleNotificationsProcessor struct {
		orm     *ent.Client
		numDays int
	}
)

func NewDeleteStaleNotificationsProcessor(orm *ent.Client, numDays int) *DeleteStaleNotificationsProcessor {
	return &DeleteStaleNotificationsProcessor{
		orm:     orm,
		numDays: numDays,
	}
}
func (d *DeleteStaleNotificationsProcessor) ProcessTask(
	ctx context.Context, t *asynq.Task,
) error {

	_, err := d.orm.Notification.
		Delete().
		Where(
			notification.CreatedAtLT(time.Now().Add(time.Hour * -24 * time.Duration(d.numDays))),
		).
		Exec(ctx)

	if err != nil {
		return err
	}


	return nil
}
