package main

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/squzy/squzy/apps/agent_client/application"
	"github.com/squzy/squzy/apps/agent_client/config"
	_ "github.com/squzy/squzy/apps/agent_client/version"
	"github.com/squzy/squzy/internal/agent"
	agent_executor "github.com/squzy/squzy/internal/agent-executor"
	"github.com/squzy/squzy/internal/logger"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"os"
)

func main() {
	cfg := config.New()
	executor, err := agent_executor.New(
		agent.New(
			cpu.Percent,
			mem.SwapMemory,
			mem.VirtualMemory,
			disk.Partitions,
			disk.Usage,
			net.IOCounters,
			timestamp.Now,
		),
		cfg.GetInterval(),
	)
	if err != nil {
		logger.Fatal(err.Error())
	}
	a := application.New(
		executor,
		cfg,
		host.Info,
		application.NewStream,
		make(chan os.Signal, 1),
	)
	err = a.Run()
	if err != nil {
		logger.Fatal(err.Error())
	}
}
