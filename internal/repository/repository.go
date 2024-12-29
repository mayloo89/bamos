package repository

import "github.com/mayloo89/bamos/internal/model"

type DBRepository interface {
	Users() bool
	InsertRoute(route model.Route) error
}
