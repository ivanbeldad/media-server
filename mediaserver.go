package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

type environment struct {
	username string
	basePort string
	uid      string
	gid      string
	storage  string
}

type command interface {
	execute() error
}

type start struct{}
type stop struct{}

func (s start) execute() error {
	username, basePort, storage, err := loadStartArgs()
	if err != nil {
		return err
	}
	uid, gid, err := getUser()
	if err != nil {
		return err
	}
	env := environment{username, basePort, uid, gid, storage}
	err = composeUp(env)
	if err != nil {
		fmt.Printf("error executing docker-compose up: %s", err.Error())
		fmt.Printf("\ntrying to revert...\n")
		composeDown()
		return err
	}
	err = changeOwner(env)
	if err != nil {
		return fmt.Errorf("impossible to change owner of storage folder")
	}
	return nil
}

func (s stop) execute() error {
	username, err := loadStopArgs()
	if err != nil {
		return err
	}
	uid, gid, err := getUser()
	if err != nil {
		return err
	}
	cmd := exec.Command("docker-compose", "down")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", username))
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEDIA_UID=%s", uid))
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEDIA_GID=%s", gid))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	err := checkSudo()
	if err != nil {
		log.Fatal(err)
	}
	cmd, err := getCommand()
	if err != nil {
		printHelp(os.Stdout)
		os.Exit(0)
	}
	err = cmd.execute()
	if err != nil {
		log.Fatal(err)
	}
}

func getCommand() (command, error) {
	if len(os.Args) <= 1 {
		return nil, fmt.Errorf("not command")
	}
	if os.Args[1] != "start" && os.Args[1] != "stop" {
		return nil, fmt.Errorf("%s is not a valid command", os.Args[1])
	}
	if os.Args[1] == "start" {
		return start{}, nil
	} else {
		return stop{}, nil
	}
}

func changeOwner(e environment) error {
	cmd := exec.Command(
		"chown",
		fmt.Sprintf("%s:%s", e.uid, e.gid),
		fmt.Sprintf("%s/%s", e.storage, e.username),
		"-R",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func composeUp(e environment) error {
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", e.username))
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEDIA_USERNAME=%s", e.username))
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEDIA_BASE_PORT=%s", e.basePort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEDIA_UID=%s", e.uid))
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEDIA_GID=%s", e.gid))
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEDIA_STORAGE=%s", e.storage))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func composeDown() error {
	cmd := exec.Command("docker-compose", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing docker-compose down: %s", err.Error())
	}
	return nil
}

func checkSudo() error {
	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot take user from environment")
	}
	if u.Uid != "0" {
		return fmt.Errorf("this script must be executed with sudo")
	}
	return nil
}

func getUser() (uid string, gid string, err error) {
	uid = os.Getenv("SUDO_UID")
	gid = os.Getenv("SUDO_GID")
	if uid == "0" {
		return "", "", fmt.Errorf("this script cannot be executed as root")
	}
	return
}

func loadStartArgs() (username string, basePort string, storage string, err error) {
	args := os.Args
	if len(args) <= 4 {
		printHelp(os.Stdout)
		os.Exit(0)
	}
	username = args[2]
	basePort = args[3]
	storage = args[4]
	port, err := strconv.Atoi(basePort)
	if err != nil {
		return "", "", "", fmt.Errorf("basePort must be a number, %s is not", basePort)
	}
	if port < 10 || port > 99 {
		return "", "", "", fmt.Errorf("basePort must be a number between 10 and 99")
	}
	if storage[0] != '/' {
		return "", "", "", fmt.Errorf("storage cannot be a relative path")
	}
	if storage[len(storage)-1] == '/' {
		storage = storage[:len(storage)-1]
	}
	return
}

func loadStopArgs() (username string, err error) {
	if len(os.Args) <= 2 {
		return "", fmt.Errorf("no user to stop")
	}
	return os.Args[2], nil
}

func printHelp(w io.Writer) {
	text := fmt.Sprintf(`
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

`)
	w.Write([]byte(text))
}
