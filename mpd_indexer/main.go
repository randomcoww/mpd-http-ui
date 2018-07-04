package main

import (
	"flag"
)

var (
	logFile = flag.String("logfile", "", "MPD log file path")
	mpdUrl = flag.String("mpdurl", "localhost:6600", "MPD URL")
	esUrl = flag.String("esurl", "http://localhost:9200", "Elasticsearch URL")
)

func main() {
	flag.Parse()

	err := NewDataFeeder(*logFile, *mpdUrl, *esUrl)

	if err != nil {
		panic(err)
	}
}
