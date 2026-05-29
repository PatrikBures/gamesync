package main

import (
	"os"
	"runtime/pprof"
	"runtime/trace"
)

func main() {
	// these debug parts should be removed or made into flags in the future
	fcpu, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(fcpu)
	defer pprof.StopCPUProfile()

	ftrace, _ := os.Create("trace.out")
	trace.Start(ftrace)
	defer trace.Stop()

	Execute()
}
