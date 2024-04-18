package main

import (
	"context"
	"fmt"
	"log"
	"os"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/core/containers"
	"github.com/containerd/containerd/v2/pkg/cio"
	"github.com/containerd/containerd/v2/pkg/namespaces"
	"github.com/containerd/containerd/v2/pkg/oci"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	containerdSocket := homedir + "/" + ".containerd.sock"

	targetPlatform := "io.containerd.runc.v2"

	if err := startContainer(context.Background(), containerdSocket, targetPlatform); err != nil {
		log.Fatalf("Fatal error: %s", err)
	}

}

func startContainer(ctx context.Context, containerdSocket, targetPlatform string) error {
	client, err := containerd.New(containerdSocket, containerd.WithDefaultRuntime(targetPlatform))
	if err != nil {
		return fmt.Errorf("containerd New: %s", err)
	}
	defer client.Close()

	containerCtx := namespaces.WithNamespace(ctx, "default")

	noCheckPlatform := func(ctx context.Context, c *containerd.UnpackConfig) error {
		c.CheckPlatformSupported = false
		return nil
	}

	containerImage, err := client.Pull(
		containerCtx,
		"docker.io/library/nginx:latest",
		containerd.WithPullUnpack,
		containerd.WithPlatform("linux/arm64"),
		containerd.WithUnpackOpts([]containerd.UnpackOpt{noCheckPlatform}),
	)
	if err != nil {
		return fmt.Errorf("pull error: %s", err)
	}

	containerSpec, err := oci.GenerateSpecWithPlatform(containerCtx, client, "linux", &containers.Container{ID: "nginx"}, oci.WithImageConfig(containerImage))
	if err != nil {
		return fmt.Errorf("containerspec error: %s", err)
	}

	containerSpec.Linux.CgroupsPath = ""

	nginxContainer, err := client.NewContainer(containerCtx, "nginx1",
		containerd.WithNewSnapshot("nginx-rootfs", containerImage),
		containerd.WithSpec(containerSpec),
	)

	if err != nil {
		return fmt.Errorf("new container error: %s", err)
	}
	defer nginxContainer.Delete(containerCtx, containerd.WithSnapshotCleanup)

	task, err := nginxContainer.NewTask(containerCtx, cio.NewCreator())
	if err != nil {
		return err
	}
	defer task.Delete(containerCtx)

	pid := task.Pid()

	fmt.Printf("Container has process id: %d\n", pid)

	err = task.Start(containerCtx)
	if err != nil {
		return fmt.Errorf("start error: %s", err)
	}
	status, err := task.Wait(containerCtx)
	if err != nil {
		return fmt.Errorf("task wait error: %s", err)
	}

	exitStatus := <-status
	fmt.Printf("Task exited with status: %v\n", exitStatus.ExitCode())

	return nil
}
