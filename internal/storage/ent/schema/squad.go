package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Squad holds the schema definition for the Squad entity.
type Squad struct {
	ent.Schema
}

// Fields of the Squad.
func (Squad) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Time("createdAt").
			Default(time.Now),
	}
}

// Edges of the Squad.
func (Squad) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("members", User.Type).
			Ref("squads"),
	}
}
