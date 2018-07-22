// https://pragmacoders.com/building-a-json-api-in-golang/
// https://gowebexamples.com/routes-using-gorilla-mux/
// https://hakaselogs.me/2017-06-23/rest-api-with-golang

package main

import (
	"flag"
)

var (
	listenUrl = flag.String("listenurl", "localhost:3000", "Listen URL")
	mpdUrl = flag.String("mpdurl", "localhost:6600", "MPD URL")
	esUrl = flag.String("esurl", "http://localhost:9200", "Elasticsearch URL")
)

func main() {
	flag.Parse()

	NewServer(*listenUrl, *mpdUrl, *esUrl)
}
