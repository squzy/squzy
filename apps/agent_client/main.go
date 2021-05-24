package main

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os"
	"squzy/apps/agent_client/application"
	"squzy/apps/agent_client/config"
	_ "squzy/apps/agent_client/version"
	"squzy/internal/agent"
	agent_executor "squzy/internal/agent-executor"
	"squzy/internal/logger"
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
			ptypes.TimestampNow,
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
