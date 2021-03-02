package main1

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	conString      = "mongodb://localhost:27017"
	dbName         = "test"
	collectionName = "railway1"
	limit          = 10
)

//var wg sync.WaitGroup
var (
	ch = make(chan int, limit)
)

type Data struct {
	TrainNo   string `bson:"trainNo"`
	TrainName string `bson:"trainName"`
	SEQ       string `bson:"seq"`
	Code      string `bson:"code"`
	StName    string `bson:"stName"`
	ATime     string `bson:"aTime"`
	DTime     string `bson:"dTime"`
	Distance  string `bson:"distance"`
	SS        string `bson:"ss"`
	SsName    string `bson:"ssName"`
	DS        string `bson:"ds"`
	DsName    string `bson:"dsName"`
}

func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func dbConn() (*mongo.Collection, *mongo.Client) {
	clientOptions := options.Client().ApplyURI(conString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	fmt.Println("Connected to MongoDB!")

	return collection, client

}
func getallTrains(w http.ResponseWriter, r *http.Request) {

	collection, client := dbConn()

	defer client.Disconnect(context.TODO())

	cursor, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}
	var trains []Data
	if err = cursor.All(context.TODO(), &trains); err != nil {
		log.Fatal(err)
	}
	bytedata, err := json.MarshalIndent(trains, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytedata)
}

func findTrains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	st := r.URL.Query().Get("stname")
	// dtn := r.URL.Query().Get("dtname")
	tno := r.URL.Query().Get("trainno")
	tname := r.URL.Query().Get("trainname")
	at := r.URL.Query().Get("arrivaltime")
	dt := r.URL.Query().Get("departuretime")

	filter := bson.D{}

	if st != "" {
		filter = append(filter, bson.E{"stName", st})
	}

	// if dtn != "" {
	// 	filter = append(filter, bson.E{"dtName", dtn})
	// }

	if tno != "" {
		filter = append(filter, bson.E{"trainNo", tno})
	}

	if tname != "" {
		filter = append(filter, bson.E{"trainName", tname})
	}

	if at != "" {
		filter = append(filter, bson.E{"aTime", at})
	}

	if dt != "" {
		filter = append(filter, bson.E{"dTime", dt})
	}

	fmt.Println(filter)
	collection, client := dbConn()
	// filter := bson.D{{"stName": st}} //bson.D{{}} specifies 'all documents'
	// filter := bson.M{"stName": st}
	issues := []Data{}
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		fmt.Println(findError)
	}
	defer client.Disconnect(context.TODO())
	for cur.Next(context.TODO()) {
		t := Data{}
		err := cur.Decode(&t)
		if err != nil {
			fmt.Println(err)
		}
		issues = append(issues, t)
	}
	defer cur.Close(context.TODO())
	if len(issues) == 0 {
		fmt.Println(mongo.ErrNoDocuments)
	}
	// fmt.Println(issues)
	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
		Data    []Data
	}{false, "Data Fetched Successfully", issues})
	w.Write(res)
}

func findOnRouteTrains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	st := r.URL.Query().Get("stname")
	dtn := r.URL.Query().Get("dtname")
	fmt.Println(st, dtn)
	// filter := bson.D{}

	// if st != "" {
	// 	filter = append(filter, bson.E{"stName", st})
	// }

	// if dtn != "" {
	// 	filter = append(filter, bson.E{"trainData.stName", dtn})
	// }
	// filter = append(filter, bson.E{"gt", "$seq < $trainData.seq"})

	collection, client := dbConn()
	defer client.Disconnect(context.TODO())

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "railway1"}, {"localField", "trainNo"}, {"foreignField", "trainNo"}, {"as", "trainData"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$trainData"}, {"preserveNullAndEmptyArrays", false}}}}
	filterStage := bson.D{{"$match", bson.D{{"stName", st}, {"trainData.stName", dtn}}}}

	fmt.Println("hi")
	// issues := []Data{}

	// query := bson.D{{
	// 	"$lookup", bson.M{
	// 		"from":         "railway1",
	// 		"localField":   "_id",
	// 		"foreignField": "_id",
	// 		"as":           "railwaydata"}}}

	// pipe := collection.Pipe(query)
	// err := pipe.All(&issues)

	// showLoadedCursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{query})
	showLoadedCursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unwindStage, filterStage})

	var showsLoaded []bson.M

	if err = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		panic(err)
	}

	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
		Data    []bson.M
	}{false, "Data Fetched Successfully", showsLoaded})
	w.Write(res)
}

func insertData() {

	collection, client := dbConn()

	defer client.Disconnect(context.TODO())

	csvData, err := ReadCsv("Indian_railway1.csv")
	if err != nil {
		fmt.Println(err)
	}

	csvData = csvData[1:50]

	for _, line := range csvData {
		ch <- 1
		func() {
			data := Data{
				TrainNo:   line[0],
				TrainName: line[1],
				SEQ:       line[2],
				Code:      line[3],
				StName:    line[4],
				ATime:     line[5],
				DTime:     line[6],
				Distance:  line[7],
				SS:        line[8],
				SsName:    line[9],
				DS:        line[10],
				DsName:    line[11],
			}

			_, err := collection.InsertOne(context.TODO(), data)
			if err != nil {
				fmt.Println(err)
			}
			<-ch
		}()
	}
	for i := 0; i < limit; i++ {
		ch <- 1
	}
}

func main() {
	start := time.Now()
	read := flag.Bool("insert", false, "a bool")
	flag.Parse()
	if *read {
		insertData()
	} else {
		fmt.Println("failed")
	}

	elp := time.Since(start)
	fmt.Println(elp)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/Trains", getallTrains)
	http.HandleFunc("/find", findTrains)
	http.HandleFunc("/route", findOnRouteTrains)
	http.ListenAndServe(":8080", nil)
}
