package pprofutil

import (
	"go.uber.org/zap"
	"os"
	"runtime"
	"runtime/pprof"
	"shop/logger"
)

func StartCpuProf() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		logger.Error.Println("create cpu profile file error: ", zap.Error(err))
		return
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		logger.Error.Println("can not start cpu profile,  error: ", zap.Error(err))
		f.Close()
	}
}

func StopCpuProf() {
	pprof.StopCPUProfile()
}

//--------Mem
func ProfGc() {
	runtime.GC() // get up-to-date statistics
}

func SaveMemProf() {
	f, err := os.Create("mem.prof")
	if err != nil {
		logger.Error.Println("create mem profile file error: ", zap.Error(err))
		return
	}

	if err := pprof.WriteHeapProfile(f); err != nil {
		logger.Error.Println("could not write memory profile: ", zap.Error(err))
	}
	f.Close()
}

// goroutine block
func SaveBlockProfile() {
	f, err := os.Create("block.prof")
	if err != nil {
		logger.Error.Println("create mem profile file error: ", zap.Error(err))
		return
	}

	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		logger.Error.Println("could not write block profile: ", zap.Error(err))
	}
	f.Close()
}
