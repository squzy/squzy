package main

import (
	"fmt"
	"squzy/apps/internal/agent"
)

func main() {
	a := agent.New()
	fmt.Println(a.GetStat())
}