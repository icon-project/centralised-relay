package dockerutil

import (
	"context"
	"fmt"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"net"
	"strings"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

type ContainerLifecycle struct {
	log               *zap.Logger
	client            *dockerclient.Client
	containerName     string
	id                string
	preStartListeners Listeners
}

func NewContainerLifecycle(log *zap.Logger, client *dockerclient.Client, containerName string) *ContainerLifecycle {
	return &ContainerLifecycle{
		log:           log,
		client:        client,
		containerName: containerName,
	}
}

func (c *ContainerLifecycle) CreateContainer(
	ctx context.Context,
	testName string,
	networkID string,
	image ibc.DockerImage,
	ports nat.PortSet,
	volumeBinds []string,
	hostName string,
	cmd []string,
) error {
	imageRef := image.Ref()
	c.log.Info(
		"Will run command",
		zap.String("image", imageRef),
		zap.String("container", c.containerName),
		zap.String("command", strings.Join(cmd, " ")),
	)

	pb, listeners, err := GeneratePortBindings(ports)
	if err != nil {
		return fmt.Errorf("failed to generate port bindings: %w", err)
	}

	c.preStartListeners = listeners

	cc, err := c.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageRef,

			Entrypoint: []string{},
			Cmd:        cmd,

			Hostname: hostName,

			Labels: map[string]string{CleanupLabel: testName},

			ExposedPorts: ports,
		},
		&container.HostConfig{
			Binds:           volumeBinds,
			PortBindings:    pb,
			PublishAllPorts: true,
			AutoRemove:      false,
			DNS:             []string{},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				networkID: {},
			},
		},
		c.containerName,
	)
	if err != nil {
		listeners.CloseAll()
		c.preStartListeners = []net.Listener{}
		return err
	}
	c.id = cc.ID
	return nil
}

func (c *ContainerLifecycle) StartContainer(ctx context.Context) error {
	// lock port allocation for the time between freeing the ports from the
	// temporary listeners to the consumption of the ports by the container
	mu.RLock()
	defer mu.RUnlock()

	c.preStartListeners.CloseAll()
	c.preStartListeners = []net.Listener{}

	if err := StartContainer(ctx, c.client, c.id); err != nil {
		return err
	}

	c.log.Info("Container started", zap.String("container", c.containerName))

	return nil
}

func (c *ContainerLifecycle) StopContainer(ctx context.Context) error {
	duration := time.Second * 30
	return c.client.ContainerStop(ctx, c.id, &duration)
}

func (c *ContainerLifecycle) RemoveContainer(ctx context.Context) error {
	err := c.client.ContainerRemove(ctx, c.id, dockertypes.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
	if err != nil && !errdefs.IsNotFound(err) {
		return fmt.Errorf("remove container %s: %w", c.containerName, err)
	}
	return nil
}

func (c *ContainerLifecycle) ContainerID() string {
	return c.id
}

func (c *ContainerLifecycle) GetHostPorts(ctx context.Context, portIDs ...string) ([]string, error) {
	cjson, err := c.client.ContainerInspect(ctx, c.id)
	if err != nil {
		return nil, err
	}
	ports := make([]string, len(portIDs))
	for i, p := range portIDs {
		ports[i] = GetHostPort(cjson, p)
	}
	return ports, nil
}
