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

type PollRepository struct {
	coll *mongo.Collection
}

var ErrPollAlreadyExist = errors.New("can't create poll: already exists")

func NewPollRepository(c *mongo.Client) *PollRepository {

	coll := c.Database("poll-service").Collection("polls")

	return &PollRepository{coll: coll}
}

func (p *PollRepository) CreatePoll(ctx context.Context, poll entities.CreatePollRequest) error {
	const op = `repository.mongo.CreatePoll`

	var existingPoll entities.CreatePollRequest

	filter := bson.M{"poll_id": poll.PollID}

	err := p.coll.FindOne(ctx, filter).Decode(&existingPoll)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return ErrPollAlreadyExist
	}

	_, err = p.coll.InsertOne(ctx, poll)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (p *PollRepository) Poll(ctx context.Context, poll entities.PollRequest) error {
	const op = `repository.mongo.Poll`

	filter := bson.M{"poll_id": poll.PollID, "poll.variants.variant_id": poll.VariantID}
	update := bson.M{"$inc": bson.M{
		"poll.variants.$.votes": 1,
	}}

	res, err := p.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("%s %w", op, repoerrs.ErrPollNotExist)
	}

	return nil
}

func (p *PollRepository) Result(ctx context.Context, poll entities.GetResultRequest) (entities.Poll, error) {
	const op = `repository.mongo.GetResult`

	var pollResult entities.CreatePollRequest

	filter := bson.M{"poll_id": poll.PollID}

	res := p.coll.FindOne(ctx, filter)

	if err := res.Decode(&pollResult); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Poll{}, fmt.Errorf("%s %w", op, repoerrs.ErrPollNotExist)
		}
		return entities.Poll{}, fmt.Errorf("%s %w", op, err)
	}
	fmt.Println(res.Raw())
	return pollResult.Poll, nil

}

func (p *PollRepository) DeletePoll(ctx context.Context, poll entities.DeleteRequest) error {
	const op = `repository.mongo.DeletePoll`

	filter := bson.M{"poll_id": poll.PollID}

	_, err := p.coll.DeleteOne(ctx, filter)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}
