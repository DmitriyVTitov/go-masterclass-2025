package memdb

import (
	"context"
	"sync"
	"ugc/internal/db"
	"ugc/internal/types"

	"go.opentelemetry.io/otel/trace"
)

type MemDB struct {
	mu   sync.Mutex
	data []types.Review
}

func New() db.DB {
	return &MemDB{}
}

func (db *MemDB) ObjectReviews(ctx context.Context, objectID int) ([]types.Review, error) {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "DB:Get")
	defer span.End()

	db.mu.Lock()
	defer db.mu.Unlock()

	var reviews []types.Review
	for _, review := range db.data {
		if review.ObjectID == objectID {
			reviews = append(reviews, review)
		}
	}

	return reviews, nil
}

func (db *MemDB) AddReview(ctx context.Context, review types.Review) error {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "DB:Set")
	defer span.End()

	db.mu.Lock()
	defer db.mu.Unlock()

	db.data = append(db.data, review)

	return nil
}
