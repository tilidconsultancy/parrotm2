package adapters

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(opts *options.ClientOptions) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	return client
}

func NewClientOptions(cfg *viper.Viper) *options.ClientOptions {
	cs := cfg.GetString("mongodb.connectionString")
	return options.Client().ApplyURI(cs)
}

func NewMongoDatabase(
	cli *mongo.Client,
	cfg *viper.Viper) *mongo.Database {

	ndb := cfg.GetString("mongodb.database")
	db := cli.Database(ndb)
	return db
}
