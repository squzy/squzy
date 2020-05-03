package main

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"squzy/apps/agent_client/application"
	"squzy/apps/agent_client/config"
	_ "squzy/apps/agent_client/version"
	"squzy/internal/agent"
	agent_executor "squzy/internal/agent-executor"
	"squzy/internal/grpcTools"
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
