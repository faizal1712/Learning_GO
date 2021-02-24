package main1

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CsvRow struct {
	TrainNo       string `bson:"TrainNo"`
	TrainName     string `bson:"TrainName"`
	StartingPoint string `bson:"StartingPoint"`
	EndingPoint   string `bson:"EndingPoint"`
}

const (
	CONNECTIONSTRING = "mongodb+srv://root:1712@cluster0.ynb7n.mongodb.net/test"
	DB               = "indian_railway"
	COLLECTION       = "indian_railway"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

func GetMongoClient() (*mongo.Client, error) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

func main() {

	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://root:1712@cluster0.ynb7n.mongodb.net/test"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Disconnect(ctx)

	start := time.Now()

	// insertRows()

	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", fs)
	http.HandleFunc("/fetch_data", FetchData)

	http.ListenAndServe(":8080", nil)
}

func insertRows() {
	client, err := GetMongoClient()
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	// Create a handle to the respective collection in the database.
	collection := client.Database(DB).Collection(COLLECTION)

	lines, err := ReadCsv("All_Indian_Trains.csv")
	if err != nil {
		panic(err)
	}

	lines = lines[1:]

	// Loop through lines & turn into object
	for _, line := range lines {
		data := CsvRow{
			TrainNo:       line[1],
			TrainName:     line[2],
			StartingPoint: line[3],
			EndingPoint:   line[4],
		}
		_, err = collection.InsertOne(context.TODO(), data)
		if err != nil {
			panic(err)
		}
		// fmt.Println(data)
	}
}

func FetchData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	client, err := GetMongoClient()
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	collection := client.Database(DB).Collection(COLLECTION)

	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	issues := []CsvRow{}
	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		panic(findError)
	}
	//Map result to slice
	for cur.Next(context.TODO()) {
		t := CsvRow{}
		err := cur.Decode(&t)
		if err != nil {
			panic(err)
		}
		issues = append(issues, t)
	}
	// once exhausted, close the cursor
	cur.Close(context.TODO())
	if len(issues) == 0 {
		panic(mongo.ErrNoDocuments)
	}
	// fmt.Println(issues)

	res, _ := json.Marshal(issues)
	w.Write(res)

}

func ReadCsv(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
