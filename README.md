### Backend setup

    docker-compose up --build

Elasticsearch:
- http://localhost:9200

MPD (control):
- tcp: localhost:6600

MPD (stream):
- http://localhost:8000

REST API:
- http://localhost:3000

Remove containers to test reindexing:

    docker-compose rm -f

### Start frontend

    cd ui
    npm run dev

#### Environment test

Connect to control using `ncmpcpp`:

    ncmpcpp -h localhost -p 6600

Naviagte client to start playing one of the test audio files. Turn on repeat mode (r) to make sure something is playing.

Open http://localhost:8000/mpd in browser to stream.
