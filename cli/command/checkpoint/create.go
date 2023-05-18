package checkpoint

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

type createOptions struct {
	container     string
	checkpoint    string
	checkpointDir string
	parentPath    string
	leaveRunning  bool
	preDump       bool
	tcpConnect    bool
	shellJob      bool
	pageServer    string
}

func newCreateCommand(dockerCli command.Cli) *cobra.Command {
	var opts createOptions

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] CONTAINER CHECKPOINT",
		Short: "Create a checkpoint from a running container",
		Args:  cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.container = args[0]
			opts.checkpoint = args[1]
			return runCreate(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.leaveRunning, "leave-running", false, "Leave the container running after checkpoint")
	flags.BoolVar(&opts.preDump, "pre-dump", false, "Pre-dump is used to pre-copy live migration")
	flags.StringVarP(&opts.parentPath, "parent-path", "", "", "Parent-Path is the path of last iteration dump image")
	flags.StringVarP(&opts.checkpointDir, "checkpoint-dir", "", "", "Use a custom checkpoint storage directory")
	flags.StringVarP(&opts.pageServer, "page-server", "", "", "Page-server is the IP:Port of page server")
	flags.BoolVar(&opts.tcpConnect, "tcp-established", false, "tcp-established is used to tcp-established live migration")
	flags.BoolVar(&opts.shellJob, "shell-job", false, "Shell-job is used to migrate terminal live migration")
	return cmd
}

func runCreate(dockerCli command.Cli, opts createOptions) error {
	client := dockerCli.Client()

	checkpointOpts := types.CheckpointCreateOptions{
		CheckpointID:  opts.checkpoint,
		CheckpointDir: opts.checkpointDir,
		PreDump:       opts.preDump,
		ParentPath:    opts.parentPath,
		Exit:          !opts.leaveRunning,
		TcpConnect:    opts.tcpConnect,
		ShellJob:      opts.shellJob,
		PageServer:    opts.pageServer,
	}

	err := client.CheckpointCreate(context.Background(), opts.container, checkpointOpts)
	if err != nil {
		return err
	}

	fmt.Fprintf(dockerCli.Out(), "%s\n", opts.checkpoint)
	return nil
}
