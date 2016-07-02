package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-respawn"
	app.Usage = "Restart Docker containers that fail health-check"
	app.UsageText = "docker-respawn <image name>"
	app.Version = "0.1"
	app.Action = func(c *cli.Context) error {
		imageName := c.Args().Get(0)
		if imageName == "" {
			cli.ShowAppHelp(c)
			return nil
		}
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		fmt.Println("Connecting to docker.sock. Checking heatlh status of", imageName)
		cli, err := client.NewClient("unix:///var/run/docker.sock", "1.24", nil, defaultHeaders)
		if err != nil {
			return err
		}
		args := filters.NewArgs()
		args.Add("Type", events.DaemonEventType)
		options := types.EventsOptions{Filters: args}
		readCloser, err := cli.Events(context.Background(), options)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(readCloser)
		for scanner.Scan() {
			var event events.Message
			err = json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&event)
			if err != nil {
				return err
			}
			action := event.Action
			actor := event.Actor
			fmt.Println(actor)
			fmt.Println(action)
		}
		return nil
	}
	app.Run(os.Args)
}
