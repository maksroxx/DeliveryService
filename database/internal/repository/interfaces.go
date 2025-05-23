package repository

import (
	"context"

	"github.com/maksroxx/DeliveryService/database/internal/models"
)

type RouteRepository interface {
	GetByID(ctx context.Context, id string) (*models.Package, error)
	GetAllRoutes(ctx context.Context, filter models.PackageFilter) ([]*models.Package, error)
	Create(ctx context.Context, route *models.Package) (*models.Package, error)
	UpdateRoute(ctx context.Context, id string, update models.PackageUpdate) (*models.Package, error)
	DeleteRoute(ctx context.Context, id string) error
	Ping(ctx context.Context) error
}
