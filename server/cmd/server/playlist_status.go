package server

import (
	"strconv"
	// "fmt"
)

type PlaylistStatus struct {
	version int
	length  int
	// playlist last updated index
	lastChangePosStart int
	lastChangePosEnd   int
}

func NewPlaylistStatus() *PlaylistStatus {
	p := &PlaylistStatus{}
	// p.update()

	// attrs, _ := mpdClient.Conn.PlaylistInfo(p.lastChangePosStart, p.lastChangePosEnd)
	// fmt.Printf("playlist: %v", attrs)

	return p
}

func (p *PlaylistStatus) update() {
	attrs, err := mpdClient.Conn.Status()
	if err != nil {
		return
	}

	version, err := strconv.Atoi(attrs["playlist"])
	if err == nil {
		if p.version < version {
			// Update change positions if version updated
			p.updatePlaylistChangePos()
			p.version = version
		}
	}

	length, err := strconv.Atoi(attrs["playlistlength"])
	if err != nil {
		p.length = length
	}
}

func (p *PlaylistStatus) updatePlaylistChangePos() {
	attrs, err := mpdClient.PlChangePosId(p.version, -1, -1)
	if err != nil {
		return
	}

	if len(attrs) > 0 {
		v, ok := attrs[0]["cpos"]
		if ok {
			i, err := strconv.Atoi(v)
			if err == nil {
				p.lastChangePosStart = i
			}
		}

		v, ok = attrs[len(attrs)-1]["cpos"]
		if ok {
			i, err := strconv.Atoi(v)
			if err == nil {
				p.lastChangePosEnd = i
			}
		}
	}
}

// keep playlist in memory so that changes can be better reported
// func (p *PlaylistStatus) syncPlaylist(batchSize int) {

// 	start := p.lastChangePosStart
// 	end := start + batchSize

// 	for end < p.lastChangePosEnd {
// 		attrs, err := mpdClient.Conn.PlaylistInfo(start, end)

// 		fmt.Prinf("playlist: %v", attrs)
// 	}
// }
