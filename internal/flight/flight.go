package flight

import "time"

type Flight struct {
	ID                 int
	Date               time.Time
	TakeoffLocation    string
	IgcFilePath        string
	UserID             int
	GliderID           int
	FlightStatisticsID int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
