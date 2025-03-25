package db

import (
	"context"

	"ugc/internal/types"
)

type DB interface {
	ObjectReviews(ctx context.Context, objectID int) ([]types.Review, error)
	AddReview(ctx context.Context, review types.Review) error
}
