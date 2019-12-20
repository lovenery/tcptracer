package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"github.com/lovenery/tcptracer/pkg/tcptracer"
)

func main() {
	fmt.Println("[main] Call tcptracer.IsStopped()")
	for {
		if tcptracer.IsStopped() {
			break
		}
	}

	fmt.Println("[main] Call tcptracer.Start()")
	tcptracer.Start()

	fmt.Println("[main] Call tcptracer.IsReady()")
	for {
		if tcptracer.IsReady() {
			break
		}
	}

	fmt.Println("[main] Starting curl google.com once")
	for i := 0; i < 1; i++ {
		cmd := exec.Command("curl", "google.com")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("[main][%d] Command finished with error: %v", i, err)
		}
	}

	fmt.Printf("[main] Call tcptracer.Stop()\n")
	tcptracer.Stop()

	fmt.Println("[main] Press CTRL+C to quit")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	fmt.Println("[main] Bye")
}