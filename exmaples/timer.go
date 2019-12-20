package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
	"github.com/lovenery/tcptracer/pkg/tcptracer"
)

func main() {
	fmt.Println("[main] Call tcptracer.Start()")
	tcptracer.Start()

	for i := 10; i > 0; i-- {
		fmt.Printf("[main] Call tcptracer.Stop() in %d seconds\n", i)
		time.Sleep(time.Duration(1) * time.Second)
	}
	fmt.Println("[main] Call tcptracer.Stop()")
	tcptracer.Stop()

	fmt.Println("[main] Press CTRL+C to quit")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	fmt.Println("[main] Bye")
}