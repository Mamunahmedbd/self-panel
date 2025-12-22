package types

import (
	"time"

	"github.com/mikestefanello/pagoda/ent"
)

type ISPProfileData struct {
	Client          *ent.ClientUser
	CurrentPackage  *ent.PackagePlan
	Usage           ISPUsageStats
	Payments        []*ent.ClientTxn
	Sessions        []*ent.RadAcct
	Tickets         []*ent.Ticket
	Balance         float64
	AvailableBalance float64
	ValidUntil      *time.Time
	AutoRenew       bool
	PackageStatus   string // "Active", "Expired", etc.
}

type ISPUsageStats struct {
	Today   uint64 // in bytes
	Weekly  uint64
	Monthly uint64
	Total   uint64
}

type TicketForm struct {
	Subject     string `form:"subject" validate:"required"`
	Description string `form:"description" validate:"required"`
	Submission  FormSubmission
}
