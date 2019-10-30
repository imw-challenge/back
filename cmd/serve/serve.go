package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/imw-challenge/back/db"
)

var (
	dataPath  string
	batchSize int
)

func main() {
	flag.StringVar(&dataPath, "datapath", "./data.csv", "path to file containing csv message data")
	flag.IntVar(&batchSize, "batchsize", 100, "maximum transaction batch size for adding messages to databse")
	flag.Parse()

	mdb, err := db.InitMessageDB()
	if err != nil {
		log.Fatal(err)
	}

	mdb.LoadFromCSV(dataPath, batchSize)
	msgs, err := mdb.FetchAntiChrono()
	if err != nil {
		panic(err)
	}

	for _, m := range msgs {
		messageLocation := time.FixedZone("", m.TZ)
		zuluTime := time.Unix(m.Time, 0)
		fmt.Printf("message from %s at %s\n", m.Name, zuluTime.In(messageLocation).Format(time.RFC3339))
	}

}

//instantiate api and register routes
//listen
