#### Backend setup

    docker-compose up --build

Elasticsearch:
- http://localhost:9200

MPD (control):
- tcp: localhost:6600

MPD (stream):
- http://localhost:8000

REST API:
- http://localhost:3000

ES data remains on the container and won't be rebuilt each run. Remove containers to force rebuild:

    docker-compose rm -f

#### UI setup

    cd mpd_ui
    npm install
    npm run dev
    
using Docker

    cd mpd_ui
    docker run -it --rm --net host -p 8080:8080 -v `pwd`:/src --entrypoint=/bin/sh node:alpine
    cd /src
    npm install
    npm run dev

#### Environment test

Connect to control using `ncmpcpp`:

    ncmpcpp -h localhost -p 6600

Naviagte client to start playing one of the test audio files. Turn on repeat mode (r) to make sure something is playing.

Open http://localhost:8000/mpd in browser to stream.
