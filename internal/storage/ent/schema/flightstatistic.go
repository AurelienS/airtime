package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FlightStatistic holds the schema definition for the FlightStatistic entity.
type FlightStatistic struct {
	ent.Schema
}

// Fields of the FlightStatistic.
func (FlightStatistic) Fields() []ent.Field {
	return []ent.Field{
		field.Int("totalThermicTime"),
		field.Int("totalFlightTime"),
		field.Int("maxClimb"),
		field.Float("maxClimbRate"),
		field.Int("totalClimb"),
		field.Float("averageClimbRate"),
		field.Int("numberOfThermals"),
		field.Float("percentageThermic"),
		field.Int("maxAltitude"),
	}
}

// Edges of the FlightStatistic.
func (FlightStatistic) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("flight", Flight.Type).
			Ref("statistic").
			Unique(),
	}
}
