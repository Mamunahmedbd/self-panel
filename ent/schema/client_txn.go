package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ClientTxn holds the schema definition for the ClientTxn entity.
type ClientTxn struct {
	ent.Schema
}

// Annotations of the ClientTxn.
func (ClientTxn) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "client_txn"},
	}
}

// Fields of the ClientTxn.
func (ClientTxn) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique(),
		field.String("transaction_ref").
			Unique().
			MaxLen(64),
		field.Float("amount").
			StorageKey("balance"). // The database column is 'balance' but it represents the transaction amount
			Default(0.00),
		field.Enum("type").
			Values("ACTIVE", "RENEWAL", "REFUND", "TRANSFER_REFUND", "TRANSFER_RECEIVED", "AUTO_RENEWAL", "PACKAGE_MIGRATION", "ADVANCE_PAYMENT"),
		field.Enum("status").
			Values("pending", "completed", "failed", "reversed").
			Default("completed"),
		field.Float("total_balance").
			Default(0.00),
		field.Enum("payment_method").
			Values("vendor_balance", "client_balance", "cash", "bank_transfer", "mobile_banking", "card", "gateway_sslcommerz", "gateway_bkash", "gateway_nagad", "gateway_stripe", "gateway_paypal", "free", "other").
			Optional().
			Nillable(),
		field.String("client_username").
			Optional().
			MaxLen(255),
		field.String("description").
			Optional(),
		field.Time("transaction_date").
			Default(time.Now),
		field.String("created_by").
			MaxLen(255),
	}
}

// Indexes of the ClientTxn.
func (ClientTxn) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("transaction_ref"),
		index.Fields("client_username"),
		index.Fields("status"),
		index.Fields("type"),
		index.Fields("transaction_date"),
	}
}

// Edges of the ClientTxn.
func (ClientTxn) Edges() []ent.Edge {
	return nil
}
