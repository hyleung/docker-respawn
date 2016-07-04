package main

import (
	"bufio"
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "respawn ", log.LstdFlags)
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
		logger.Printf("Connecting to docker.sock. Checking heatlh status of", imageName)
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
			if actor.Attributes["image"] == imageName {
				logger.Println(actor)
				logger.Println(action)
				if action == "health_status: unhealthy" {
					//get the continer so that we can re-start it
					args := filters.NewArgs()
					args.Add("id", actor.ID)
					listOptions := types.ContainerListOptions{Filter: args}
					containers, err := cli.ContainerList(context.Background(), listOptions)
					if err != nil {
						logger.Println(err)
						return err
					}
					log.Println(containers[0])
					logger.Println("Stopping", actor.Attributes["name"], "due to failed health check")
					timeout := 10 * time.Second
					stopErr := cli.ContainerStop(context.Background(), actor.ID, &timeout)
					if stopErr != nil {
						logger.Println("Failed to stop container", actor.ID)
						return err
					}
				}
			}
		}
		return nil
	}
	app.Run(os.Args)
}
