package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PackagePlan holds the schema definition for the PackagePlan entity.
type PackagePlan struct {
	ent.Schema
}

// Annotations of the PackagePlan.
func (PackagePlan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "packages"},
	}
}

// Fields of the PackagePlan.
func (PackagePlan) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Unique(),
		field.String("name").
			NotEmpty().
			MaxLen(255),
		field.String("pool_name").
			Optional().
			MaxLen(255),
		field.String("profile_name").
			NotEmpty().
			MaxLen(100),
		field.Float("price").
			Default(0.00),
		field.String("currency").
			Default("BDT").
			MaxLen(3),
		field.Bool("is_active").
			Default(true),
		field.Time("created_date").
			Default(time.Now).
			Immutable(),
	}
}

// Indexes of the PackagePlan.
func (PackagePlan) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("profile_name"),
		index.Fields("is_active"),
	}
}

// Edges of the PackagePlan.
func (PackagePlan) Edges() []ent.Edge {
	return nil
}
