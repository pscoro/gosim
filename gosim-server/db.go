package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DATABASE_NAME = "gosim"
const USER_COLLECTION = "users"

type dbConn struct {
	client *mongo.Client
	db     *mongo.Database
}

func newDatabaseConnection() (*dbConn, error) {
	godotenv.Load()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("URI")).
		SetServerAPIOptions(serverAPIOptions)
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to database", err)
		return nil, err
	}
	db := client.Database(DATABASE_NAME)

	return &dbConn{client: client, db: db}, nil
}

func (dbc *dbConn) registerUser(u *user) error {
	coll := dbc.db.Collection(USER_COLLECTION)
	type userSchema struct {
		name  string
		email string
	}
	us := userSchema{
		name:  u.name,
		email: u.email,
	}
	insertRes, err := coll.InsertOne(context.TODO(), us)
	if err != nil {
		log.Println("Error inserting user in db", err)
		return err
	}
	t, ok := insertRes.InsertedID.(uint32)
	if ok == false {
		log.Println("InsertedID can not convert to uint32", err)
		return err
	}
	u.id = t
	return nil
}

func (dbc *dbConn) closeDatabaseConnection() {
	err := dbc.client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
}
