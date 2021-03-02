package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

type PageId struct {
	ID int `json:"id"`
}

type CsvRow struct {
	TrainNo       string `bson:"TrainNo"`
	TrainName     string `bson:"TrainName"`
	StartingPoint string `bson:"StartingPoint"`
	EndingPoint   string `bson:"EndingPoint"`
}

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load("env_var.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func GetMongoClient() (*mongo.Client, error) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		CONNECTIONSTRING := goDotEnvVariable("CONNECTIONSTRING")
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
	fmt.Println("hello")

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

	LIMIT, _ := strconv.Atoi(goDotEnvVariable("LIMIT"))

	ch := make(chan bool, LIMIT)
	start := time.Now()
	fgptr := flag.Bool("insert", false, "a bool")
	flag.Parse()

	if *fgptr {
		insertRows(ch)
	} else {
		fmt.Println(*fgptr)
	}

	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", fs)
	http.HandleFunc("/fetch_data", FetchData)

	http.ListenAndServe(":8080", nil)
}

func insertRows(ch chan bool) {
	client, err := GetMongoClient()
	defer client.Disconnect(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
	}

	DB := goDotEnvVariable("DB")
	COLLECTION := goDotEnvVariable("COLLECTION")

	// Create a handle to the respective collection in the database.
	collection := client.Database(DB).Collection(COLLECTION)

	lines, err := ReadCsv("All_Indian_Tasks.csv")
	if err != nil {
		fmt.Println(err)
	}

	lines = lines[1:]

	// Loop through lines & turn into object
	for _, line := range lines {
		ch <- true
		go func(line []string, ch chan bool) {
			data := CsvRow{
				TrainNo:       line[1],
				TrainName:     line[2],
				StartingPoint: line[3],
				EndingPoint:   line[4],
			}
			_, err = collection.InsertOne(context.TODO(), data)
			if err != nil {
				fmt.Println(err)
			}
			<-ch
		}(line, ch)
		// fmt.Printf("var1 = %T\n", line)
		// fmt.Println(data)
	}
	LIMIT, _ := strconv.Atoi(goDotEnvVariable("LIMIT"))

	for i := 0; i < LIMIT; i++ {
		ch <- true
	}
}

func FetchData(w http.ResponseWriter, r *http.Request) {
	// id, _ := strconv.Atoi(r.FormValue("id"))
	bodydata, _ := ioutil.ReadAll(r.Body)
	// fmt.Println(bodydata)
	var score PageId
	_ = json.Unmarshal(bodydata, &score)
	fmt.Println(score)

	id := score.ID

	w.Header().Set("Content-Type", "application/json")
	client, err := GetMongoClient()
	if err != nil {
		fmt.Println(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
	}

	DB := goDotEnvVariable("DB")
	COLLECTION := goDotEnvVariable("COLLECTION")

	collection := client.Database(DB).Collection(COLLECTION)

	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	issues := []CsvRow{}

	option := options.Find()
	option.SetLimit(10)
	sk := 10 * id
	fmt.Println(sk)
	option.SetSkip(int64(sk))

	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter, option)
	if findError != nil {
		panic(findError)
	}
	//Map result to slice
	for cur.Next(context.TODO()) {
		t := CsvRow{}
		err := cur.Decode(&t)
		if err != nil {
			fmt.Println(err)
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
