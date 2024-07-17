package mongo

import (
	"context"
	"errors"
	"fmt"
	"poll-service/internal/entities"
	"poll-service/internal/repository/repoerrs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserVotesRepository struct {
	coll *mongo.Collection
}

func NewUserVotesRepository(c *mongo.Client) *UserVotesRepository {
	return &UserVotesRepository{
		coll: c.Database("poll-service").Collection("requests"),
	}
}

func (r *UserVotesRepository) Save(ctx context.Context, req entities.Request) error {
	op := `repository.mongo.request.Save`

	filter := bson.M{"remote_addr": req.RemoteAddr, "request.pollid": req.Request.PollID}

	var existingRequest entities.Request

	err := r.coll.FindOne(ctx, filter).Decode(&existingRequest)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("%s %w", op, repoerrs.ErrUserAlreadyVoted)
	}

	_, err = r.coll.InsertOne(ctx, req)
	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}

	return nil
}

func (r *UserVotesRepository) ClearByPollID(ctx context.Context, poll entities.DeleteRequest) error {
	op := `repository.mongo.request.ClearByPollID`

	if err := r.coll.Drop(ctx); err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}
