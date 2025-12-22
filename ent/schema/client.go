package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ClientUser holds the schema definition for the ClientUser entity (PPPoE users).
type ClientUser struct {
	ent.Schema
}

// Annotations of the ClientUser.
func (ClientUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "clients"},
	}
}

// Fields of the ClientUser.
func (ClientUser) Fields() []ent.Field {
	return []ent.Field{
		// Primary Key
		field.Int("id").
			Positive().
			Unique(),

		// User Information
		field.String("name").
			NotEmpty().
			MaxLen(255),
		field.String("username").
			NotEmpty().
			MaxLen(255).
			Unique(),
		field.String("password").
			NotEmpty().
			MaxLen(255).
			Sensitive(), // Plain-text password for PPPoE
		field.String("mobile_number").
			NotEmpty().
			MaxLen(255),
		field.String("email").
			NotEmpty().
			MaxLen(255),
		field.String("photo").
			Optional().
			MaxLen(255),
		field.Text("description").
			Optional(),

		// Financial
		field.Float("balance").
			Default(0.00).
			Min(0),

		// Address Information
		field.String("address_line1").
			Optional().
			MaxLen(255),
		field.String("address_line2").
			Optional().
			MaxLen(255),
		field.String("city").
			Optional().
			MaxLen(255),
		field.String("district").
			Optional().
			MaxLen(255),
		field.String("upazila").
			Optional().
			MaxLen(255),
		field.String("union_name").
			Optional().
			MaxLen(255),
		field.String("zip").
			Optional().
			MaxLen(255),

		// Status and Payment
		field.Enum("status").
			Values("active", "inactive").
			Default("inactive"),
		field.Time("payment_date").
			Optional().
			Nillable(),
		field.String("payment_type").
			Optional().
			MaxLen(255),
		field.Bool("auto_renew").
			Default(true),

		// Company and Vendor Association
		field.String("c_name").
			NotEmpty().
			MaxLen(32),
		field.Int("vendor_id").
			Positive(),

		// Package Assignment
		field.String("package_pool").
			Optional().
			MaxLen(255),
		field.String("user_profile").
			Optional().
			MaxLen(64).
			Comment("Current user profile (from packages.profile_name)"),
		field.String("next_user_profile").
			Optional().
			MaxLen(255).
			Comment("Next profile for renewal (from packages.profile_name)"),

		// Audit Fields
		field.String("created_by").
			NotEmpty().
			MaxLen(255),
		field.String("updated_by").
			Optional().
			MaxLen(255),
		field.Time("created_date").
			Default(time.Now).
			Immutable(),
		field.Time("updated_date").
			Optional().
			Nillable().
			UpdateDefault(time.Now),
	}
}

// Indexes of the ClientUser.
func (ClientUser) Indexes() []ent.Index {
	return []ent.Index{
		// Single column indexes
		index.Fields("mobile_number"),
		index.Fields("email"),
		index.Fields("c_name"),
		index.Fields("vendor_id"),
		index.Fields("package_pool"),
		index.Fields("user_profile"),
		index.Fields("next_user_profile"),
		index.Fields("auto_renew"),
		index.Fields("created_date"),
		index.Fields("created_by"),
		index.Fields("status"),
		index.Fields("name"),
		index.Fields("balance"),

		// Composite indexes for performance
		index.Fields("c_name", "status", "created_date"),
		index.Fields("c_name", "vendor_id", "status", "created_date"),
		index.Fields("package_pool", "status", "created_date"),
		index.Fields("auto_renew", "status", "c_name", "created_date"),
		index.Fields("vendor_id", "package_pool", "status"),
		index.Fields("c_name", "vendor_id", "status"),
		index.Fields("updated_date", "status"),
	}
}

// Edges of the ClientUser.
func (ClientUser) Edges() []ent.Edge {
	return nil
}
