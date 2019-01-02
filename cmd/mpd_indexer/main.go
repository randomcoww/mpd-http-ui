package mpd_indexer

import (
	"flag"

	es "github.com/randomcoww/go-mpd-es/pkg/es_handler"
	mpd "github.com/randomcoww/go-mpd-es/pkg/mpd_handler"
	"github.com/sirupsen/logrus"
)

func Main() {
	logFile := flag.String("logfile", "", "MPD log file path")
	mpdURL := flag.String("mpdurl", "localhost:6600", "MPD URL")
	esURL := flag.String("esurl", "http://localhost:9200", "Elasticsearch URL")
	flag.Parse()

	newDataFeeder(*logFile, *mpdURL, *esURL)
}

// process to read log to create add and remove events
func newDataFeeder(logFile, mpdURL, esURL string) {
	logrus.Infof("Create MPD log pipe: %s", logFile)

	e := newLogEvents()
	go e.readLog(logFile)

	// Start services
	mpdClient := mpd.NewMpdClient("tcp", mpdURL)
	esClient := es.NewEsClient(esURL, esIndex, esDocument, esMapping)

	<-mpdClient.Ready
	<-esClient.Ready

	// Start parseing logged events after MPD is accessible
	e.parseLog(mpdClient, esClient)
}