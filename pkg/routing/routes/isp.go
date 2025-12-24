package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/repos/msg"
	"github.com/mikestefanello/pagoda/pkg/types"
)

type ispRoutes struct {
	ctr controller.Controller
}

func NewISPRoutes(ctr controller.Controller) *ispRoutes {
	return &ispRoutes{
		ctr: ctr,
	}
}

func (c *ispRoutes) GetTickets(ctx echo.Context) error {
	// For now, redirect back to profile where tickets are listed
	// or implement a dedicated tickets page if needed later
	return c.ctr.Redirect(ctx, "profile")
}

func (c *ispRoutes) CreateTicket(ctx echo.Context) error {
	var form types.TicketForm
	if err := ctx.Bind(&form); err != nil {
		return err
	}

	if err := form.Submission.Process(ctx, &form); err != nil {
		return err
	}

	if form.Submission.HasErrors() {
		return c.ctr.Redirect(ctx, "profile")
	}

	// Save ticket to DB
	client, err := c.ctr.Container.GetAuthenticatedClient(ctx)
	if err != nil {
		return err
	}

	_, err = c.ctr.Container.ORM.Ticket.
		Create().
		SetClientID(client.ID).
		SetClientUsername(client.Username).
		SetSubject(form.Subject).
		SetDescription(form.Description).
		SetStatus("open").
		SetPriority("medium").
		Save(ctx.Request().Context())

	if err != nil {
		return c.ctr.Fail(err, "failed to create ticket")
	}

	msg.Success(ctx, "Support ticket opened successfully. Our team will review it soon.")
	return c.ctr.Redirect(ctx, "profile")
}

func (c *ispRoutes) AddFunds(ctx echo.Context) error {
	// Implement balance load logic
	return c.ctr.Redirect(ctx, "profile")
}

func (c *ispRoutes) RenewPackage(ctx echo.Context) error {
	// Implement package renewal logic
	return c.ctr.Redirect(ctx, "profile")
}

func (c *ispRoutes) ChangePlan(ctx echo.Context) error {
	// Implement plan change logic
	return c.ctr.Redirect(ctx, "profile")
}

func (c *ispRoutes) ToggleAutoRenew(ctx echo.Context) error {
	client, err := c.ctr.Container.GetAuthenticatedClient(ctx)
	if err != nil {
		return err
	}

	// Toggle the current state
	newStatus := !client.AutoRenew

	_, err = c.ctr.Container.ORM.ClientUser.
		UpdateOne(client).
		SetAutoRenew(newStatus).
		Save(ctx.Request().Context())

	if err != nil {
		return c.ctr.Fail(err, "failed to update settings")
	}

	if newStatus {
		msg.Success(ctx, "Auto-renew enabled")
	} else {
		msg.Info(ctx, "Auto-renew disabled")
	}

	return c.ctr.Redirect(ctx, "profile")
}
