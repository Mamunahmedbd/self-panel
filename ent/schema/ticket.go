package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Ticket holds the schema definition for the Ticket entity.
type Ticket struct {
	ent.Schema
}

// Fields of the Ticket.
func (Ticket) Fields() []ent.Field {
	return []ent.Field{
		field.String("subject").
			NotEmpty().
			MaxLen(255),
		field.Text("description").
			NotEmpty(),
		field.Enum("status").
			Values("open", "pending", "closed").
			Default("open"),
		field.Enum("priority").
			Values("low", "medium", "high").
			Default("medium"),
		field.Int("client_id").
			Positive(),
		field.String("client_username").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes of the Ticket.
func (Ticket) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("client_id"),
		index.Fields("client_username"),
		index.Fields("status"),
		index.Fields("created_at"),
	}
}

// Edges of the Ticket.
func (Ticket) Edges() []ent.Edge {
	return nil
}
