package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Flight holds the schema definition for the Flight entity.
type Flight struct {
	ent.Schema
}

// Fields of the Flight.
func (Flight) Fields() []ent.Field {
	return []ent.Field{
		field.Time("date"),
		field.String("takeoffLocation"),
		field.String("igcFilePath"),
	}
}

// Edges of the Flight.
func (Flight) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pilot", User.Type).
			Ref("flights").
			Unique(),
		edge.To("statistic", FlightStatistic.Type).
			Unique(),
	}
}
