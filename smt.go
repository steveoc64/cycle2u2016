///////////////////////////////////////////////////////////////////////////////////////////
//  SMT - We have multi cores, lets use them !

package main

import (
	"fmt"
	"runtime"
)

func _initSMT() {

	numCores := runtime.NumCPU()

	useCores := numCores / 2
	if useCores < 1 {
		useCores = 1
	}
	fmt.Println("Cycle2U server: - running on", useCores, "of", numCores, "CPU Cores")

	runtime.GOMAXPROCS(useCores)
}
