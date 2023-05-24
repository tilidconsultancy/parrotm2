package adapters

import (
	"context"
	"log"
	"pm2/internal/ports"

	"reflect"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbRepository[T interface{}] struct {
	collection *mongo.Collection
}

func NewMongoDbRepository[T interface{}](
	db *mongo.Database) ports.Repository[T] {

	var r T
	coll := db.Collection(strings.ToLower(reflect.TypeOf(r).Name()))
	return &MongoDbRepository[T]{
		collection: coll,
	}
}

func (r *MongoDbRepository[T]) GetAll(
	ctx context.Context,
	filter map[string]interface{}) []T {

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	result := []T{}
	for cur.Next(ctx) {
		var el T
		err = cur.Decode(&el)
		if err != nil {
			panic(err)
		}
		result = append(result, el)
	}

	return result
}

func (r *MongoDbRepository[T]) GetAllSkipTake(
	ctx context.Context,
	filter map[string]interface{},
	skip int64,
	take int64) *ports.Pagination[T] {
	ct, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		panic(err)
	}

	op := options.Find()
	op.SetSkip(skip)
	op.SetSort(map[string]interface{}{
		"_id": -1,
	})
	if take > 0 {
		op.SetLimit(take)
	}
	cur, err := r.collection.Find(ctx, filter, op)

	if err != nil {
		panic(err)
	}
	result := []T{}
	for cur.Next(ctx) {
		var el T
		err = cur.Decode(&el)
		if err != nil {
			panic(err)
		}
		result = append(result, el)
	}

	return &ports.Pagination[T]{
		Data:  result,
		Count: ct,
	}
}

func (r *MongoDbRepository[T]) Count(ctx context.Context,
	filter map[string]interface{}) int64 {
	ct, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		panic(err)
	}
	return ct
}

func (r *MongoDbRepository[T]) GetFirst(
	ctx context.Context,
	filter map[string]interface{}) *T {
	var el T
	err := r.collection.FindOne(ctx, filter).Decode(&el)

	if err == mongo.ErrNoDocuments {
		return nil
	}

	if err != nil {
		panic(err)
	}

	return &el
}

func (r *MongoDbRepository[T]) Insert(
	ctx context.Context,
	entity *T) {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		panic(err)
	}
}

func (r *MongoDbRepository[T]) InsertAll(
	ctx context.Context,
	entities []T) {

	var uis []interface{}
	for _, ui := range entities {
		uis = append(uis, ui)
	}
	_, err := r.collection.InsertMany(ctx, uis)
	if err != nil {
		panic(err)
	}
}

func (r *MongoDbRepository[T]) Replace(
	ctx context.Context,
	filter map[string]interface{},
	entity *T) {
	up := true
	_, err := r.collection.ReplaceOne(ctx, filter, entity, &options.ReplaceOptions{
		Upsert: &up,
	})
	if err != nil {
		panic(err)
	}
}

func (r *MongoDbRepository[T]) DeleteAll(
	ctx context.Context,
	filter map[string]interface{}) {
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		panic(err)
	}
}

func (r *MongoDbRepository[T]) Aggregate(
	ctx context.Context,
	pipeline []map[string]interface{}) []map[string]interface{} {
	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		panic(err)
	}
	result := []map[string]interface{}{}
	for cur.Next(ctx) {
		var el map[string]interface{}
		err = cur.Decode(&el)
		if err != nil {
			panic(err)
		}
		result = append(result, parseMongoUUIDs(el))
	}

	return result
}

func parseMongoUUIDs(parentMap map[string]interface{}) map[string]interface{} {
	for key, value := range parentMap {
		if reflect.TypeOf(value) == reflect.TypeOf(map[string]interface{}{}) {
			parentMap[key] = parseMongoUUIDs(value.(map[string]interface{}))
			continue
		}
		if reflect.TypeOf(value) == reflect.TypeOf(primitive.Binary{}) {
			if id, err := uuid.FromBytes(value.(primitive.Binary).Data); err == nil {
				parentMap[key] = id
			} else {
				log.Print(err.Error())
			}
		}
	}
	return parentMap
}
