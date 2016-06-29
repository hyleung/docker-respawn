package main

import (
	"fmt"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func main() {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	fmt.Println("Connecting to docker.sock")
	cli, err := client.NewClient("unix:///var/run/docker.sock", "1.24", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}
	options := types.ContainerListOptions{All: false}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		fmt.Println(c.Image, c.Names)
	}
}
