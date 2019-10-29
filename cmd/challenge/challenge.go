package main

import (
	"flag"
	"log"
)

var (
	dataPath string
)

flag.StringVar(&dataPath, "datapath", "./data.csv", "path to file containing csv message data")
flag.Parse()


//Create memdb and load csv
//instantiate api and register routes
//listen
