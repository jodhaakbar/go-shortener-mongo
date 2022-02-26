package storemongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

type Url struct {
	Shorturl    string `bson:"shorturl"`
	Originalurl string `bson:"originalurl"`
	Uuid        string `bson:"uuid"`
	Webhook     string `bson:"webhook"`
	Createdat   int64  `bson:"createdat"`
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func connect() (*mongo.Database, error) {

	// godotenv package
	mongourl := goDotEnvVariable("MONGODB_URL")

	fmt.Printf("godotenv : %s = %s \n", "STRONGEST_AVENGER", mongourl)

	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://" + mongourl)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("goshort"), nil
}

func SaveUrlMapping(shortUrl string, originalUrl string, UUID string, Webhook string) {

	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	// check if exist
	count, err := db.Collection("urlCollection").CountDocuments(ctx, bson.M{"shorturl": shortUrl})
	if err != nil {
		panic(err)
	}
	if count == 0 {
		_, err = db.Collection("urlCollection").InsertOne(ctx, Url{Shorturl: shortUrl, Originalurl: originalUrl, Uuid: UUID, Createdat: time.Now().Unix(), Webhook: Webhook})

		if err != nil {
			panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
		}
	}

	fmt.Printf("Saved shortUrl: %s - originalUrl: %s\n", shortUrl, originalUrl)
}

func RetrieveInitialUrl(shortUrl string) string {
	var resurl Url

	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.Collection("urlCollection").FindOne(ctx, bson.M{"shorturl": shortUrl}).Decode(&resurl); err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, shortUrl))
	}

	res := make([]Url, 0)
	res = append(res, resurl)

	return res[0].Originalurl
}
