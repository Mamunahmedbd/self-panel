package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RadAcct holds the schema definition for the RadAcct entity.
type RadAcct struct {
	ent.Schema
}

// Annotations of the RadAcct.
func (RadAcct) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "radacct"},
	}
}

// Fields of the RadAcct.
func (RadAcct) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			StorageKey("radacctid").
			Unique(),
		field.String("acctsessionid").
			MaxLen(64),
		field.String("acctuniqueid").
			Unique().
			MaxLen(32),
		field.String("username").
			MaxLen(64),
		field.Time("acctstarttime").
			Optional().
			Nillable(),
		field.Time("acctstoptime").
			Optional().
			Nillable(),
		field.Uint32("acctsessiontime").
			Optional().
			Nillable(),
		field.Int64("acctinputoctets").
			Optional().
			Nillable(),
		field.Int64("acctoutputoctets").
			Optional().
			Nillable(),
		field.String("framedipaddress").
			MaxLen(15),
		field.String("acctterminatecause").
			MaxLen(32),
	}
}

// Indexes of the RadAcct.
func (RadAcct) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username"),
		index.Fields("acctstarttime"),
		index.Fields("acctstoptime"),
		index.Fields("framedipaddress"),
	}
}

// Edges of the RadAcct.
func (RadAcct) Edges() []ent.Edge {
	return nil
}
