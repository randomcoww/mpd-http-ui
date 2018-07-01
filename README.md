### Setup

    cd env
    docker-compose up --build

Elasticsearch:
- http://localhost:9200

MPD (control):
- tcp: localhost:6600

MPD (stream):
- http://localhost:8000


Need to clean containers after each run to test indexing new items each time.

    docker-compose rm -f

#### Environment test

Connect to control using `ncmpcpp`:

Run

    ncmpcpp -h localhost -p 6600

Naviagte client to start playing one of the test audio files. Turn on repeat mode (r) to make sure something is playing.

Open http://localhost:8000/mpd in browser to stream.
