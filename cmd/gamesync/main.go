package main

import (
	"log"
	"os"
	"runtime"
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

	fmem, _ := os.Create("mem.prof")
	defer fmem.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(fmem); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
