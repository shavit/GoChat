# Go Chat Server

[![Build Status](https://travis-ci.org/shavit/gochat.svg?branch=master)](https://travis-ci.org/shavit/gochat)

TCP based chat server with Docker build.

![Go TCP chat server](https://github.com/shavit/GoChat/blob/master/doc/preview.png?raw=true)

## Quick start
1. Build the app
````
docker build -t gochat .
````

2. Test
````
docker run --rm -t gochat go test ./...
````

3. Start the server.
````
docker network create gochat
docker run --rm -ti --name gochat_server --net gochat gochat
go build src/github.com/shavit/gochat
./gochat
````

In the `Dockerfile` the __$SERVER\_HOSTNAME__ is set to __gochat\_server__, so make sure to name the server __gochat\_server__.

4. Connect to the server
````
docker run --rm -ti --net gochat gochat
./gochat
````

### Development
````
docker run --rm -ti --net gochat -v $PWD:/app gochat
````
