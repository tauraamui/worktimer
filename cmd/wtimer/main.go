package main

import (
	"flag"
	"time"

	"github.com/tauraamui/worktimer/internal/tui"
)

var (
	workDuration  = flag.Duration("wd", 45*time.Minute, "work duration")
	breakDuration = flag.Duration("bd", 15*time.Minute, "break duration")
)

func main() {
	flag.Parse()
	tui.StartTea(*workDuration, *breakDuration)
}
