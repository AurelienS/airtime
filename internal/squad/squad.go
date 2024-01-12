package squad

import (
	"time"
)

type Squad struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Members   []Member
}

type Member struct {
	ID       int
	SquadID  int
	UserID   int
	Admin    bool
	JoinedAt time.Time
}
