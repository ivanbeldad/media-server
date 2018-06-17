# Media Server

Run a set of containers in one step. Included plex server, sonarr, radarr,
jackett, transmission, ddclient and nginx.

## Requirements

- Your preferred flavour of linux
- Golang
- Docker CE
- Git?
- A cup of coffee

## Motivation

Create and manage **independent** media centers scoped by user.

## Where to start

Simply clone this repository and execute the mediaserver script
```
git clone https://github.com/ivandelabeldad/media-server.git
cd media-server
sudo go run mediaserver.go start awesome 42 /storage
```

### Help

```
Start o stop the desired user media server

Commands:
	start username basePort storage
	stop username

Arguments:
	username		Name of the user owner of the media server
	basePort		Number between 10 and 99 used as port prefix
	storage			Path where all will be stored

Example:
	sudo mediaserver ivandelabeldad 42 /media

	Storage:
		/media/ivandelabeldad
	Ports:
		Plex					4200
		Tranmission		4201
		Sonarr				4202
		Radarr				4203
		Jackett				4204
```

## Important

Plex doesn't work well inside a bridge network, so in order to use it and prevent
potential problems a host or macvlan network is recommended.

Due to the impossibility of create multiple plex instances using the host network, the
obviuous option is macvlan.

A basic example (mostly home networks could use it using 192.168.X.0/24)
```
docker network create -d macvlan --gateway=192.168.1.1 --subnet=192.168.1.0/24 -o parent=eth0 --ip-range=192.168.1.16/28 macvlan
```

## License

Media Server is open-sourced software licensed under
the [MIT license](https://github.com/ivandelabeldad/media-server/blob/master/LICENSE).
