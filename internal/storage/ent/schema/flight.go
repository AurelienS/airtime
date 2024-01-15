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
		field.String("takeoffLocation"),
		field.String("igcFilePath"),
	}
}

func (Flight) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pilot", User.Type).Unique().
			Ref("flights"),
		edge.To("statistic", FlightStatistic.Type).
			Unique(),
	}
}

func (Flight) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("date").Edges("pilot").Unique(),
	}
}
