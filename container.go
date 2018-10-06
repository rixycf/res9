package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/client"
	"os"
)

// rescue func inspect unhealth container, and revive it.
func rescue() {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("create client error: ", err)
		os.Exit(1)
	}

	// =============================================
	listOp := types.ContainerListOptions{}
	containerList, err := cli.ContainerList(ctx, listOp)
	if err != nil {
		fmt.Println("show container list error:", err)
		os.Exit(1)
	}

	for _, cl := range containerList {
		resultInspect, err := cli.ContainerInspect(ctx, cl.ID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if getHealthStatus(resultInspect) == "unhealthy" {
			fmt.Println("unhealhy container ID: ", resultInspect.ID)
			err = reviveContainer(ctx, cli, resultInspect)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
	// ===============================================
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
