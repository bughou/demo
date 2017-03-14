package main

import (
	"fmt"

	"github.com/bughou-go/xiaomei/xiaomei/images"
	"github.com/bughou-go/xiaomei/xiaomei/stack"
	"github.com/bughou-go/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
)

// Run, Build, Push, Deploy and Ps commands
func commonCmds(svcName string) []*cobra.Command {
	cmds := []*cobra.Command{}
	if svcName != `` {
		cmds = append(cmds, runCmd(svcName))
	}
	var target, s string
	if svcName == `` {
		target, s = `all`, `s`
	} else {
		target, s = svcName, ``
	}
	cmds = append(cmds,
		buildCmd(svcName, target, s),
		pushCmd(svcName, target, s),
		deployCmd(svcName, target, s),
		psCmd(svcName, target, s),
	)
	return cmds
}

func runCmd(svcName string) *cobra.Command {
	var publish []string
	cmd := &cobra.Command{
		Use:   `run`,
		Short: fmt.Sprintf(`run    %s image.`, svcName),
		RunE: z.NoArgCall(func() error {
			return images.Run(svcName, publish)
		}),
	}
	cmd.Flags().StringSliceVarP(&publish, `publish`, `p`, nil, `publish ports for container.`)
	return cmd
}

func buildCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `build`,
		Short: fmt.Sprintf(`build  %s image%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return images.Build(svcName)
		}),
	}
}

func pushCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `push`,
		Short: fmt.Sprintf(`push   %s image%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return images.Push(svcName)
		}),
	}
}

func deployCmd(svcName, target, s string) *cobra.Command {
	var noBuild, noPush bool
	cmd := &cobra.Command{
		Use:   `deploy`,
		Short: fmt.Sprintf(`deploy %s service%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return stack.Deploy(svcName, noBuild, noPush)
		}),
	}
	cmd.Flags().BoolVarP(&noBuild, `no-build`, `B`, false, fmt.Sprintf(`do not build the image%s.`, s))
	cmd.Flags().BoolVarP(&noPush, `no-push`, `P`, false, fmt.Sprintf(`do not push the image%s.`, s))
	return cmd
}

func psCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `ps`,
		Short: fmt.Sprintf(`list tasks of %s service%s.`, target, s),
		RunE: func(c *cobra.Command, args []string) error {
			return stack.Ps(svcName, args)
		},
	}
}
