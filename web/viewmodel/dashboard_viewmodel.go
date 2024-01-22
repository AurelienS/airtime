package viewmodel

type DashboardStatsView struct {
	FlightCount       string
	TotalFlightTime   string
	TotalDistance     string
	AverageFlightTime string
	MaxFlightTime     string
	MaxDistance       string
}
type DashboardCurrentYearStatsView struct {
	FlightCount       string
	TotalFlightTime   string
	TotalDistance     string
	AverageFlightTime string
	MaxFlightTime     string
	MaxDistance       string
}

type DashboardSitesStatsView struct {
	Name            string
	FlightCount     string
	TotalFlightTime string
	TotalDistance   string
}

type DashboardFlightView struct {
	FlightNumber    string
	Date            string
	TakeoffLocation string
	TotalFlightTime string
	TotalDistance   string
	Link            string
}
type DashboardView struct {
	User            UserView
	CurrentYearStat DashboardCurrentYearStatsView
	LastFlights     []DashboardFlightView
	AllTimeStats    DashboardStatsView
	SitesStats      DashboardSitesStatsView
}
