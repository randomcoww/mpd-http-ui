package server

import (
	"flag"

	"github.com/randomcoww/go-mpd-es/pkg/elasticsearch"
	"github.com/randomcoww/go-mpd-es/pkg/mpd"
	"github.com/sirupsen/logrus"
)

var (
	listenurl = flag.String("listenurl", "", "Listen URL")
	logFile   = flag.String("logfile", "", "MPD log file path")
	mpdSocket = flag.String("mpdsocket", "/run/mpd/socket", "MPD Socket")
	esUrl     = flag.String("esurl", "http://localhost:9200", "Elasticsearch URL")
)

var (
	mpdLogReader   *MpdLogEvents
	mpdClient      *mpd.MpdClient
	mpdEvent       *mpd.MpdEvent
	esClient       *elasticsearch.EsClient
	playlistStatus *PlaylistStatus
)

func Main() {
	flag.Parse()

	var (
		err error
	)

	exit := make(chan struct{})

	mpdLogReader, err = NewMpdLogReader(*logFile)
	if err != nil {
		logrus.Errorf("Could not open MPD log, %v", err)
		panic("Could not open MPD log")
	}

	mpdClient = mpd.NewMpdClient("unix", *mpdSocket)
	mpdEvent = mpd.NewMpdEvent("unix", *mpdSocket)
	esClient = elasticsearch.NewEsClient(*esUrl, esSongIndex, esSongDocument, esSongMapping)

	// playlistStatus = NewPlaylistStatus()

	go runLogIndexer()
	go runEventHandler()

	<-exit
}

// Read mpd logs and index to ES
func runLogIndexer() {
	for {
		select {
		case e := <-mpdLogReader.AddEvent:
			logrus.Infof("Add item event: %s", e)
			attr := mpdClient.GetDatabaseItem(e)

			logrus.Infof("Add item: %v", attr)

			esClient.IndexBulk(e, Song{
				File:     e,
				Date:     attr["date"],
				Duration: attr["duration"],
				Composer: attr["composer"],
				Album:    attr["album"],
				Track:    attr["track"],
				Title:    attr["title"],
				Artist:   attr["artist"],
				Genre:    attr["genre"],
			})
		case e := <-mpdLogReader.DeleteEvent:
			logrus.Infof("Delete item event: %s", e)
			esClient.DeleteBluk(e)
		}
	}
}

// handle events from MPD
func runEventHandler() {
	for {
		select {
		case e := <-mpdEvent.Events:

			switch e {
			case "player":

			case "playlist":

			case "mixer", "options", "outputs":

			case "update":
			}
		}
	}
}
