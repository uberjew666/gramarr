# gramarr
## A [Radarr](https://github.com/Radarr/Radarr), [Sonarr](https://github.com/Sonarr/Sonarr) and [Lidarr](https://github.com/Lidarr/Lidarr) Telegram Bot featuring user authentication/level access.

## Features

### Sonarr

- Search for TV Shows by name.
- Pick which seasons you want to download.
- Choose which quality and language profile you want to download.

### Radarr

- Search for Movies by name.
- Choose which quality profile you want to download.

### Lidarr

- Search for Artists by name.
- Choose which quality profile you want to download.

---

## Requirements

- A running instance of Radarr
- A running instance of Sonarr
- A running instance of Lidarr

### If running from source

- [Go](https://golang.org/)

### If running from docker

- [Docker](https://docker.io)
- [Docker Compose](https://docs.docker.com/compose/)

---

## Configuration

- Copy the `config.json.template` file to `config.json` and set-up your configuration;

---

## Running it

### From Docker

```bash
$ docker-compose up -d
```

Alternatively:

```bash
$ docker run -d --name gramarr uberjew666/gramarr:latest
```

### From source

```bash
$ go get github.com/uberjew666/gramarr
$ cd $GOPATH/src/github.com/uberjew666/gramarr
$ go get
$ go run .
```
