package reviews

import (
	"context"
	"testing"
	"ugc/internal/db/memdb"
	"ugc/internal/types"
)

func TestService_ObjectReviews(t *testing.T) {
	ctx := context.Background()

	s := New(memdb.New())

	r := types.Review{
		ObjectID: 1,
		Text:     "text",
	}
	s.AddReview(ctx, r)
	s.AddReview(ctx, r)

	tests := []struct {
		name     string
		objectID int
		wantLen  int
		wantErr  bool
	}{
		{
			objectID: 1,
			wantLen:  2,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ObjectReviews(ctx, tt.objectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ObjectReviews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Service.ObjectReviews() = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}
