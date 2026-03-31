package db

import "github.com/thirteen37/setapp/internal/model"

// Store is the interface for database access. The concrete *DB type satisfies it.
type Store interface {
	AllApps() ([]model.App, error)
	SearchApps(query string) ([]model.App, error)
	FindApp(name string) (*model.App, error)
	LoadCategories(apps []model.App) error
	AppCategories(pk int) ([]string, error)
	AllCategories() ([]model.Category, error)
	AppsByCategory(categoryName string) ([]model.App, error)
	Close() error
}
