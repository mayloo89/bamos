package database

import (
	"context"
	"time"

	"github.com/mayloo89/bamos/internal/model"
)

func (db *postgreDatabase) Users() bool {
	return true
}

func (m *postgreDatabase) InsertRoute(route model.Route) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := `INSERT INTO routes (route_id, agency_id, route_short_name, route_long_name, route_desc, route_type, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := m.DB.ExecContext(ctx, query,
		route.RouteID,
		route.AgencyID,
		route.ShortName,
		route.LongName,
		route.Desc,
		route.Type,
		time.Now(),
		time.Now())

	if err != nil {
		return err
	}

	return nil
}
