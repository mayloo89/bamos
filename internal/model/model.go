package model

import "time"

// Route is the representation of the route table in the database.
type Route struct {
	ID        int
	RouteID   int
	AgencyID  int
	ShortName string
	LongName  string
	Desc      string
	Type      int
	CreatedAt time.Time
	UpdatedAt time.Time
}
