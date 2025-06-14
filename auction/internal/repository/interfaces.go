package repository

import (
	"context"

	"github.com/maksroxx/DeliveryService/auction/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bidder interface {
	PlaceBid(ctx context.Context, bid *models.Bid) error
	GetBidsByPackage(ctx context.Context, packageID string) ([]*models.Bid, error)
	WatchBidsByPackage(ctx context.Context, packageID string) (*mongo.ChangeStream, error)
	GetTopBidByPackage(ctx context.Context, packageID string) (*models.Bid, error)
}

type Packager interface {
	Create(ctx context.Context, pkg *models.Package) (*models.Package, error)
	Update(ctx context.Context, pkg *models.Package) error
	FindByID(ctx context.Context, packageID string) (*models.Package, error)
	FindUserPackages(ctx context.Context, userId string) ([]*models.Package, error)
	FindByFailedStatus(ctx context.Context) ([]*models.Package, error)
	FindByAuctioningStatus(ctx context.Context) ([]*models.Package, error)
	FindByWaitingStatus(ctx context.Context) ([]*models.Package, error)
}
