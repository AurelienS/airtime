package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Flight struct {
	ent.Schema
}

func (Flight) Fields() []ent.Field {
	return []ent.Field{
		field.Time("date"),
		field.String("location"),
		field.Int("duration"),
		field.Int("distance"),
		field.Int("altitudeMax"),
		field.String("igcData"),
	}
}

func (Flight) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pilot", User.Type).Unique().
			Ref("flights"),
	}
}

func (Flight) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("date").Edges("pilot").Unique(),
	}
}
