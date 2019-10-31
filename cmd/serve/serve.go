package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/imw-challenge/back/api"
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
	_, err = mdb.FetchAntiChrono()
	if err != nil {
		log.Fatal(err)
	}

	//instantiate api and register routes
	apiHandle, err := api.InitAPI(mdb)
	if err != nil {
		log.Fatal(err)
	}

	//listen
	log.Fatal(http.ListenAndServe("0.0.0.0:9000", apiHandle.GetRouter()))

}
