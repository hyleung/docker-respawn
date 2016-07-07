package main

import (
	"bufio"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/network"
	"golang.org/x/net/context"
	"os"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-respawn"
	app.Usage = "Restart Docker containers that fail health-check"
	app.UsageText = "docker-respawn <image name>"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "debug", Usage: "Enable debug level logging"},
	}
	app.Action = func(c *cli.Context) error {
		imageName := c.Args().Get(0)
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		if imageName == "" {
			cli.ShowAppHelp(c)
			return nil
		}
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		log.Info("Connecting to docker.sock... Checking heatlh status of ", imageName)
		cli, err := client.NewClient("unix:///var/run/docker.sock", "1.24", nil, defaultHeaders)
		if err != nil {
			log.Error(err)
			return err
		}
		args := filters.NewArgs()
		args.Add("Type", events.DaemonEventType)
		options := types.EventsOptions{Filters: args}
		readCloser, err := cli.Events(context.Background(), options)
		if err != nil {
			log.Error(err)
			return err
		}
		scanner := bufio.NewScanner(readCloser)
		for scanner.Scan() {
			var event events.Message
			err = json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&event)
			if err != nil {
				log.Error(err)
				return err
			}
			action := event.Action
			actor := event.Actor
			if actor.Attributes["image"] == imageName {
				log.Debug(actor)
				log.Debug(action)
				if action == "health_status: unhealthy" {
					//get the continer so that we can re-start it
					args := filters.NewArgs()
					args.Add("id", actor.ID)
					//listOptions := types.ContainerListOptions{Filter: args}
					//containers, err := cli.ContainerList(context.Background(), listOptions)
					if err != nil {
						log.Error(err)
						return err
					}
					log.Info("Stopping ", actor.Attributes["name"], " due to failed health check")
					timeout := 1 * time.Second
					stopErr := cli.ContainerStop(context.Background(), actor.ID, &timeout)
					if stopErr != nil {
						log.Error("Failed to stop container ", actor.ID)
						return err
					}
					//get the container
					stoppedContainer, err := cli.ContainerInspect(context.Background(), actor.ID)
					response, err := cli.ContainerCreate(context.Background(),
						stoppedContainer.Config,
						stoppedContainer.HostConfig,
						&network.NetworkingConfig{EndpointsConfig: stoppedContainer.NetworkSettings.Networks},
						"")
					if err != nil {
						log.Error(err)
						return err
					}
					log.Info("Created new container: ", response.ID)
					err = cli.ContainerStart(context.Background(), response.ID, types.ContainerStartOptions{})
					if err != nil {
						log.Error(err)
						return err
					}
					log.Info("Container respawned")
				}
			}
		}
		return nil
	}
	app.Run(os.Args)
}
