package dashboard

type Service struct {
	dashboardRepo Repository
}

func NewService(dashboardRepo Repository) Service {
	return Service{
		dashboardRepo: dashboardRepo,
	}
}
