package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"golang.org/x/net/context"
	"strings"
)

func main() {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	fmt.Println("Connecting to docker.sock")
	cli, err := client.NewClient("unix:///var/run/docker.sock", "1.24", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}
	options := types.EventsOptions{}
	readCloser, err := cli.Events(context.Background(), options)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		var event events.Message
		err = json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&event)
		if err != nil {
			panic(err)
		}
		fmt.Println(event)
	}
}
