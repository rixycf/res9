package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type CustomClient struct {
	client.Client
}

func NewCustomClient() (*CustomClient, error) {
	return client.NewClientWithOpts(client.FromEnv)
}

func (c *CustomClient) ReviveContainer(ctx context.Context, cj types.ContainerJSON) error {

	fmt.Println("stop container ... : ")
	err := c.ContainerStop(ctx, cj.ID, nil)
	if err != nil {
		return err
	}

	fmt.Println("remove container ...: ", cj.ID)
	err = c.ContainerRemove(ctx, cj.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	config := &container.Config{
		Image: "066e88b2a453",
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

	createbody, err := c.ContainerCreate(ctx, config, hostConfig, netconfig, "adblocker")
	if err != nil {
		return err
	}

	err = c.ContainerStart(ctx, createbody.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	fmt.Println("start container: ", createbody.ID)

	return nil
}
