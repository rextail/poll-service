package mongo

import (
	"context"
	"fmt"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	Client *mongo.Client
}

func New(connString string) (*Client, error) {
	const op = `repository.mongo.New`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, fmt.Errorf(`%s %w`, op, err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf(`%s %w`, op, err)
	}

	return &Client{Client: client}, nil
}
