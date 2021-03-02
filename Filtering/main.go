package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
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
	SEQ       int    `bson:"seq"`
	Code      string `bson:"code"`
	StName    string `bson:"stName"`
	ATime     string `bson:"aTime"`
	DTime     string `bson:"dTime"`
	Distance  int    `bson:"distance"`
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

func findRouteTrains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	st := r.URL.Query().Get("stname1")
	dtn := r.URL.Query().Get("stname2")
	fmt.Println(st, dtn)
	// filter := bson.D{}

	// if st != "" {
	// 	filter = append(filter, bson.E{"stName", st})
	// }

	// if dtn != "" {
	// 	filter = append(filter, bson.E{"trainData.stName", dtn})
	// }
	// filter = append(filter, bson.E{"trainData.seq", bson.D{{"$gt", "seq"}}})

	collection, client := dbConn()
	defer client.Disconnect(context.TODO())

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "railway1"}, {"localField", "trainNo"}, {"foreignField", "trainNo"}, {"as", "trainData"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$trainData"}, {"preserveNullAndEmptyArrays", false}}}}
	filterStage := bson.D{{"$match", bson.D{{"stName", st}, {"trainData.stName", dtn}}}}

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
		Count   int
		Data    []bson.M
	}{false, "Data Fetched Successfully", len(showsLoaded), showsLoaded})
	w.Write(res)
}

func findOneRouteTrains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	st := r.URL.Query().Get("stname")
	dtn := r.URL.Query().Get("dtname")

	collection, client := dbConn()
	defer client.Disconnect(context.TODO())

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "railway1"}, {"localField", "trainNo"}, {"foreignField", "trainNo"}, {"as", "trainData"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$trainData"}, {"preserveNullAndEmptyArrays", false}}}}
	filterStage := bson.D{{"$match", bson.D{{"stName", st}, {"trainData.stName", dtn}}}}
	// nextfilterStage := bson.D{{"$match", bson.D{{"seq", bson.D{{"$gt", "trainData.seq"}}}}}}
	fmt.Println("hi")

	showLoadedCursor, _ := collection.Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unwindStage, filterStage})

	var showsLoaded []bson.M

	count := 0
	for showLoadedCursor.Next(context.TODO()) {
		var t bson.M
		err := showLoadedCursor.Decode(&t)
		if err != nil {
			panic(err)
		}
		seq := t["seq"].(int32)
		dseq := t["trainData"].(bson.M)["seq"].(int32)

		// fmt.Println(seq, dseq)
		if seq < dseq {
			count++
			showsLoaded = append(showsLoaded, t)
		}
	}

	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
		Count   int
		Data    []bson.M
	}{false, "Data Fetched Successfully", count, showsLoaded})
	w.Write(res)
}

func findSortedOneRouteTrains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	st := r.URL.Query().Get("stname")
	dtn := r.URL.Query().Get("dtname")

	collection, client := dbConn()
	defer client.Disconnect(context.TODO())

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "railway1"}, {"localField", "trainNo"}, {"foreignField", "trainNo"}, {"as", "trainData"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$trainData"}, {"preserveNullAndEmptyArrays", false}}}}
	filterStage := bson.D{{"$match", bson.D{{"stName", st}, {"trainData.stName", dtn}}}}
	// nextfilterStage := bson.D{{"$match", bson.D{{"seq", bson.D{{"$gt", "trainData.seq"}}}}}}
	fmt.Println("hi")

	showLoadedCursor, _ := collection.Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unwindStage, filterStage})

	var showsLoaded []bson.M

	count := 0
	for showLoadedCursor.Next(context.TODO()) {
		var t bson.M
		err := showLoadedCursor.Decode(&t)
		if err != nil {
			panic(err)
		}
		seq := t["seq"].(int32)
		dseq := t["trainData"].(bson.M)["seq"].(int32)

		// fmt.Println(seq, dseq)
		if seq < dseq {
			count++

			dt := strings.ReplaceAll(t["dTime"].(string), ":", "")                       //departure time from soure
			at := strings.ReplaceAll(t["trainData"].(bson.M)["aTime"].(string), ":", "") //arrival time at destination

			dtime, err := time.Parse("150405", dt)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(dtime)

			atime, err := time.Parse("150405", at)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(atime)

			t["DeptTime"] = dtime
			t["ArriTime"] = atime

			diff := atime.Sub(dtime)
			// if(diff < 0) {
			// 	diff = dtime.Sub(atime)
			// }
			// fmt.Println(diff)
			out := time.Time{}.Add(diff).String()
			t["timeTaken"] = out[11:19]
			showsLoaded = append(showsLoaded, t)
		}
	}

	sort.Slice(showsLoaded, func(i, j int) bool {
		dtime, err := time.Parse("150405", strings.ReplaceAll(showsLoaded[i]["timeTaken"].(string), ":", ""))
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(dtime)

		atime, err := time.Parse("150405", strings.ReplaceAll(showsLoaded[j]["timeTaken"].(string), ":", ""))
		if err != nil {
			log.Fatal(err)
		}
		diff := dtime.Sub(atime)
		if diff < 0 {
			return true
		}
		return false
	})

	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
		Count   int
		Data    []bson.M
	}{false, "Data Fetched Successfully", count, showsLoaded})
	w.Write(res)
}

func findMinMaxRouteTrains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	st := r.URL.Query().Get("stname")
	dtn := r.URL.Query().Get("dtname")

	collection, client := dbConn()
	defer client.Disconnect(context.TODO())

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "railway1"}, {"localField", "trainNo"}, {"foreignField", "trainNo"}, {"as", "trainData"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$trainData"}, {"preserveNullAndEmptyArrays", false}}}}
	filterStage := bson.D{{"$match", bson.D{{"stName", st}, {"trainData.stName", dtn}}}}
	// nextfilterStage := bson.D{{"$match", bson.D{{"seq", bson.D{{"$gt", "trainData.seq"}}}}}}
	fmt.Println("hi")

	showLoadedCursor, _ := collection.Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unwindStage, filterStage})

	var showsLoaded []bson.M
	var distArray []int32

	var count = 0
	var minDist int32 = 0
	var maxDist int32 = 0
	var minDistObj bson.M
	var maxDistObj bson.M

	for showLoadedCursor.Next(context.TODO()) {
		var t bson.M
		err := showLoadedCursor.Decode(&t)
		if err != nil {
			panic(err)
		}
		// seq := t["seq"].(int32)
		// dseq := t["trainData"].(bson.M)["seq"].(int32)

		sdist := t["distance"].(int32)
		ddist := t["trainData"].(bson.M)["distance"].(int32)

		// fmt.Println(seq, dseq)
		if sdist < ddist {
			distance := ddist - sdist
			if count == 0 {
				count++
				minDist = distance
				maxDist = distance
				minDistObj = t
				maxDistObj = t
			} else {
				// fmt.Println(distance, minDist, maxDist)
				if minDist > distance {
					minDist = distance
					minDistObj = t
				} else if distance > maxDist {
					maxDist = distance
					maxDistObj = t
				}
			}
		}
	}
	count++

	showsLoaded = append(showsLoaded, minDistObj)
	showsLoaded = append(showsLoaded, maxDistObj)
	distArray = append(distArray, minDist)
	distArray = append(distArray, maxDist)

	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
		Count   int
		Dist    []int32
		Data    []bson.M
	}{false, "Data Fetched Successfully", count, distArray, showsLoaded})
	w.Write(res)
}

func insertData() {

	collection, client := dbConn()

	defer client.Disconnect(context.TODO())

	csvData, err := ReadCsv("Indian_railway1.csv")
	if err != nil {
		fmt.Println(err)
	}

	csvData = csvData[1:]

	for _, line := range csvData {
		ch <- 1
		func(line []string) {
			seq, _ := strconv.Atoi(line[2])
			dist, _ := strconv.Atoi(line[7])
			data := Data{
				TrainNo:   line[0],
				TrainName: line[1],
				SEQ:       seq,
				Code:      line[3],
				StName:    line[4],
				ATime:     line[5],
				DTime:     line[6],
				Distance:  dist,
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
		}(line)
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
	http.HandleFunc("/trains", getallTrains)
	http.HandleFunc("/find", findTrains)
	http.HandleFunc("/route", findRouteTrains)
	http.HandleFunc("/oneroute", findOneRouteTrains)
	http.HandleFunc("/sortedoneroute", findSortedOneRouteTrains)
	http.HandleFunc("/minmaxoneroute", findMinMaxRouteTrains)
	http.ListenAndServe(":8080", nil)
}
