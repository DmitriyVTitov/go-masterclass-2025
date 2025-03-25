package reviews

import (
	"context"
	"ugc/internal/db"
	"ugc/internal/types"
)

type Service struct {
	db db.DB
}

func New(db db.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) ObjectReviews(ctx context.Context, objectID int) ([]types.Review, error) {
	return s.db.ObjectReviews(ctx, objectID)
}

func (s *Service) AddReview(ctx context.Context, review types.Review) error {
	return s.db.AddReview(ctx, review)
}
