package test

import (
	"context"
	"exampleproject/db"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

func waitForDB(killFunc func() error) error {
	var err error

	for i := 0; i < 5; i++ {
		if err = db.Connect(); err == nil {
			return nil
		}

		time.Sleep(time.Second)
	}

	return err
}

func StartPostgres() (func() error, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/postgres", types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	io.Copy(os.Stdout, reader)

	containerConfig := &container.Config{
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", os.Getenv("POSTGRES_USER")),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", os.Getenv("POSTGRES_PASSWORD")),
			fmt.Sprintf("POSTGRES_DB=%s", os.Getenv("POSTGRES_DB")),
		},
		Image:        "postgres",
		ExposedPorts: nat.PortSet{"5432:5432": struct{}{}},
	}

	initSQLFilePath, err := filepath.Abs("../init.sql")
	if err != nil {
		return nil, err
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: initSQLFilePath,
				Target: "/docker-entrypoint-initdb.d/init.sql",
			},
		},
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "5432",
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	go func() {

		select {
		case err := <-errCh:
			panic(err)
		case <-statusCh:
		}

		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}()

	killFunc := func() error { return cli.ContainerKill(ctx, resp.ID, "SIGKILL") }

	if err := waitForDB(killFunc); err != nil {
		return nil, err
	}

	return killFunc, nil
}
