package viewmodel

type DashboardStatsView struct {
	FlightCount     string
	TotalDuration   string
	TotalDistance   string
	AverageDuration string

	DurationMax       string
	DurationMaxFlight FlightView

	DistanceMax       string
	DistanceMaxFlight FlightView

	AltitudeMax       string
	AltitudeMaxFlight FlightView
}

type DashboardSitesStatsView struct {
	Name            string
	FlightCount     string
	TotalFlightTime string
	TotalDistance   string
}

type DashboardView struct {
	User            UserView
	CurrentYear     string
	FirstYear       string
	LastYear        string
	LastFlights     []FlightView
	CurrentYearStat DashboardStatsView
	AllTimeStats    DashboardStatsView
	SitesStats      DashboardSitesStatsView
}
