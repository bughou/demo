package run

import (
	"fmt"
	"runtime"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/deploy"
	"github.com/lovego/xiaomei/services/images"
)

func run(env, svcName string) error {
	containerName := release.ServiceName(env, svcName) + ".run"

	args := []string{"run", "-it", "--rm", "--name=" + containerName}
	if runtime.GOOS == "linux" { // only linux support host network
		args = append(args, "--network=host")
	}
	image := images.Get(svcName)
	if portEnvVar := image.PortEnvVar(); portEnvVar != "" {
		runPort := getRunPort(image, env, svcName)
		args = append(args, "-e", fmt.Sprintf("%s=%d", portEnvVar, runPort))
		if runtime.GOOS != "linux" {
			args = append(args, fmt.Sprintf("--publish=%d:%d", runPort, runPort))
		}
	}
	if options := image.FlagsForRun(env); len(options) > 0 {
		args = append(args, options...)
	}

	args = append(args, deploy.GetCommonArgs(svcName, env, "")...)
	if err := removeContainer(containerName); err != nil {
		return err
	}
	_, err := cmd.Run(cmd.O{}, "docker", args...)
	return err
}

func getRunPort(image images.Image, env, svcName string) uint16 {
	if ports := release.GetService(env, svcName).Ports; len(ports) > 0 {
		return ports[0]
	}
	return image.DefaultPort()
}

func removeContainer(name string) error {
	if !cmd.Ok(cmd.O{}, "docker", "inspect", "-f", "{{ .State.Status }}", name) {
		return nil
	}
	_, err := cmd.Run(cmd.O{}, "docker", "rm", name)
	return err
}
