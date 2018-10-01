package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/client"
	"os"
)

func main() {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("create client error: ", err)
		os.Exit(1)
	}

	listOp := types.ContainerListOptions{}
	containerList, err := cli.ContainerList(ctx, listOp)
	if err != nil {
		fmt.Println("show container list error:", err)
		os.Exit(1)
	}

	for _, cl := range containerList {
		fmt.Println(cl.Image)

		resultInspect, err := cli.ContainerInspect(ctx, cl.ID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(resultInspect)

	}

}

func getHealthStatus(cj types.ContainerJSON) string {
	if cj.State.Health != nil {
		return cj.State.Health.Status
	}
	return ""
}

func ReviveContainer(cli *client.Client, ctx context.Context, cj types.ContainerJSON) error {

	fmt.Println("stop container ... : ")
	err := cli.ContainerStop(ctx, cj.ID, nil)
	if err != nil {
		return err
	}

	fmt.Println("remove container ...: ", cj.ID)
	err = cli.ContainerRemove(ctx, cj.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	config := &container.Config{
		Image: cj.Image,
		ExposedPorts: nat.PortSet{
			nat.Port("53/udp"): struct{}{},
		},
	}

	netconfig := &network.NetworkingConfig{}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port("53/udp"): []nat.PortBinding{{HostPort: "53"}},
		},
	}

	createbody, err := cli.ContainerCreate(ctx, config, hostConfig, netconfig, "adblocker")
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, createbody.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	fmt.Println("start container: ", createbody.ID)

	return nil
}
