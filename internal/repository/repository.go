package repository

import (
	"context"

	"poll-service/internal/entities"
	"poll-service/internal/repository/mongo"
)

type PollRepository interface {
	Poll(ctx context.Context, poll entities.PollRequest) error
	DeletePoll(ctx context.Context, poll entities.DeleteRequest) error
	CreatePoll(ctx context.Context, poll entities.CreatePollRequest) error
	Result(ctx context.Context, req entities.GetResultRequest) (entities.Poll, error)
}

type UserVotesRepository interface {
	Save(ctx context.Context, req entities.Request) error
	ClearByPollID(ctx context.Context, voteRequest entities.DeleteRequest) error
}

type Repositories struct {
	PollRepository
	UserVotesRepository
}

func New(client *mongo.Client) *Repositories {
	return &Repositories{
		PollRepository:      mongo.NewPollRepository(client.Client),
		UserVotesRepository: mongo.NewUserVotesRepository(client.Client),
	}
}
