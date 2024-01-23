package viewmodel

type DashboardStatsView struct {
	FlightCount       string
	TotalFlightTime   string
	TotalDistance     string
	AverageFlightTime string
	MaxDuration       string
	MaxDurationFlight DashboardFlightView
	MaxDistance       string
	MaxDistanceFlight DashboardFlightView
}
type DashboardCurrentYearStatsView struct {
	FlightCount       string
	TotalFlightTime   string
	TotalDistance     string
	AverageFlightTime string
	MaxDuration       string
	MaxDurationFlight DashboardFlightView
	MaxDistance       string
	MaxDistanceFlight DashboardFlightView
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
	CurrentYear     string
	FirstYear       string
	LastYear        string
	LastFlights     []DashboardFlightView
	AllTimeStats    DashboardStatsView
	SitesStats      DashboardSitesStatsView
}
