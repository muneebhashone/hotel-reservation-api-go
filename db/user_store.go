package db

import (
	"context"
	"time"

	"github.com/muneebhashone/go-fiber-api/config"
	"github.com/muneebhashone/go-fiber-api/types"
	"github.com/muneebhashone/go-fiber-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUsers(ctx context.Context, page, pageSize int, sortField, sortOrder, searchQuery string) ([]*types.User, int, error)
	GetUser(context.Context, string) (*types.User, error)
	DeleteUser(context.Context, string) error
	CreateUser(context.Context, types.User) (*types.User, error)
	UpdateUser(ctx context.Context, id string, input types.UpdateUserInput) error
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	collection := client.Database(config.DBNAME).Collection(config.USER_COLLECTION)
	return &MongoUserStore{
		client:     client,
		collection: collection,
	}
}

func (store *MongoUserStore) GetUser(ctx context.Context, identifier string) (*types.User, error) {
	var filter bson.M

	// Check if the identifier is an email
	if utils.IsEmail(identifier) {
		filter = bson.M{"email": identifier}
	} else {
		// Try to convert the identifier to an ObjectID
		objectId, err := primitive.ObjectIDFromHex(identifier)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": objectId}
	}

	var user *types.User
	err := store.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (store *MongoUserStore) CreateUser(ctx context.Context, data types.User) (*types.User, error) {
	insertedUser, err := store.collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	data.ID = insertedUser.InsertedID.(primitive.ObjectID)

	return &data, nil
}

func (store *MongoUserStore) UpdateUser(ctx context.Context, id string, input types.UpdateUserInput) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	input.UpdatedAt = time.Now()

	_, err = store.collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": input})
	if err != nil {
		return err
	}

	return nil
}

func (store *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = store.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}

	return nil
}

func (store *MongoUserStore) GetUsers(ctx context.Context, page, pageSize int, sortField, sortOrder, searchQuery string) ([]*types.User, int, error) {
	// Create the query filter for searching
	filter := bson.M{}
	if searchQuery != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"firstname": bson.M{"$regex": searchQuery, "$options": "i"}},
				{"lastname": bson.M{"$regex": searchQuery, "$options": "i"}},
				{"email": bson.M{"$regex": searchQuery, "$options": "i"}},
			},
		}
	}

	// Sorting logic
	sort := bson.M{"created_at": -1} // Default sorting by createdAt
	if sortField != "" {
		sort = bson.M{sortField: -1}
	}

	// Pagination logic
	skip := (page - 1) * pageSize

	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$facet", Value: bson.M{
			"totalData": []bson.D{{{Key: "$count", Value: "total"}}},
			"userData":  []bson.D{{{Key: "$sort", Value: sort}}, {{Key: "$skip", Value: skip}}, {{Key: "$limit", Value: pageSize}}},
		}}},
	}

	cursor, err := store.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Parsing the result
	type result struct {
		TotalData []struct{ Total int } `bson:"totalData"`
		UserData  []*types.User         `bson:"userData"`
	}
	var results []result

	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	if len(results) == 0 {
		return []*types.User{}, 0, nil
	}

	total := 0
	if len(results[0].TotalData) > 0 {
		total = results[0].TotalData[0].Total
	}

	return results[0].UserData, total, nil
}
