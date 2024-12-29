package database

import "github.com/mayloo89/bamos/internal/model"

func (db *testDatabase) Users() bool {
	return true
}

func (db *testDatabase) InsertRoute(route model.Route) error {
	return nil
}
