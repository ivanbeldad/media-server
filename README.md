# Media Server

## Usage

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
		Plex			4200
		Tranmission		4201
		Sonarr			4202
		Radarr			4203
		Jackett			4204
```

## License

Media Server is open-sourced software licensed under
the [MIT license](https://github.com/ivandelabeldad/media-server/blob/master/LICENSE).
