package main

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"squzy/apps/agent/application"
	"squzy/apps/agent/config"
	"squzy/apps/internal/agent"
	agent_executor "squzy/apps/internal/agent-executor"
	"squzy/apps/internal/grpcTools"
    _ "squzy/apps/agent/version"
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
		cfg.GetExecutionTimeout(),
	)
	if err != nil {
		log.Fatal(err)
	}
	a := application.New(
		executor,
		grpcTools.New(),
		cfg,
		host.Info,
		application.NewStream,
	)
	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
