package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/client"
)

const (
	host    = "unix:///var/run/docker.sock"
	version = "1.38"
)

// rescue func inspect unhealth container, and revive it.
func rescue() {

	ctx := context.Background()
	cli, err := client.NewClient(host, version, nil, nil)
	if err != nil {
		errlog.Println("Error: ", err)
		return
	}

	listOp := types.ContainerListOptions{}
	containerList, err := cli.ContainerList(ctx, listOp)
	if err != nil {
		errlog.Println("Error: ", err)
		return
	}

	for _, cl := range containerList {
		resultInspect, err := cli.ContainerInspect(ctx, cl.ID)
		if err != nil {
			errlog.Println("Error: ", err)
			return
		}

		if getHealthStatus(resultInspect) == "unhealthy" {
			fmt.Println("unhealhy container ID: ", resultInspect.ID)
			err = reviveContainer(ctx, cli, resultInspect)
			if err != nil {
				errlog.Println("Error: ", err)
				return
			}
		}
	}
}

func getHealthStatus(cj types.ContainerJSON) string {
	if cj.State.Health != nil {
		return cj.State.Health.Status
	}
	return ""
}

func reviveContainer(ctx context.Context, cli *client.Client, cj types.ContainerJSON) error {

	err := cli.ContainerStop(ctx, cj.ID, nil)
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, cj.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	netconfig := &network.NetworkingConfig{}
	createbody, err := cli.ContainerCreate(ctx, cj.Config, cj.HostConfig, netconfig, cj.Name)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, createbody.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}
