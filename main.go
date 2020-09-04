package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

var tstCnt = 100

func main() {
	ch1 := make(chan float64)
	ch2 := make(chan float64)
	defer close(ch1)
	defer close(ch2)

	go runRoda(ch1)
	go runHanamiApi(ch2)

	roda := <-ch1
	hanamiApi := <-ch2

	fmt.Printf("Send %v Request to Roda server and hanami-api server\n", tstCnt)
	fmt.Printf("show avarage response time below \n")
	fmt.Printf("============================\n")
	fmt.Printf("      roda    | %.3f ms\n", roda)
	fmt.Printf("   hanami-api | %.3f ms\n", hanamiApi)
	fmt.Printf("============================\n")
}

func runRoda(ch chan<- float64) {
	port := "9000"
	path := "./apps/roda/config.ru"
	pidFile := "./roda/pid"
	r := healthRack(path, port, pidFile)
	ch <- r
}

func runHanamiApi(ch chan<- float64) {
	port := "9001"
	path := "./apps/hanami-api/config.ru"
	pidFile := "./hanami-api/pid"
	r := healthRack(path, port, pidFile)
	ch <- r
}

func healthRack(path, port, pidFile string) float64 {
	ch1 := make(chan struct{}, 0)
	ch2 := make(chan float64)
	defer close(ch1)
	defer close(ch2)

	go func(_ chan<- struct{}) {
		exec.Command("rackup", path, "-p", port, "-P", pidFile).Run()
	}(ch1)

	go func(ch chan<- float64) {
		var f float64
		for i := 0; i < tstCnt; i++ {
			o, _ := exec.Command("curl", "-w", "%{time_total}", fmt.Sprintf("localhost:%v", port)).Output()
			r, _ := strconv.ParseFloat(string(o), 32)
			f += r
		}
		ch <- f / float64(tstCnt) * 1000
	}(ch2)

	v := <-ch2
	exec.Command("kill", "-9", fmt.Sprintf("`cat %v`", pidFile)).Run()
	return v
}
