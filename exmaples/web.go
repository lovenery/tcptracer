package main

import (
	"fmt"
	"net/http"
	"github.com/lovenery/tcptracer/pkg/tcptracer"
)

func TracerStart(w http.ResponseWriter, req *http.Request) {
	err := tcptracer.Start()
	if err != nil {
		fmt.Fprintf(w, "TcpTracer started to fail: %s\n", err)
	} else {
		fmt.Fprintf(w, "TcpTracer started\n")
	}
}

func TracerStop(w http.ResponseWriter, req *http.Request) {
	tcptracer.Stop()
	fmt.Fprintf(w, "TcpTracer stopped\n")
}

func TracerStopped(w http.ResponseWriter, req *http.Request) {
	status := "Not Stopped"
	if tcptracer.IsStopped() {
		status = "Stopped"
	}
	fmt.Fprintf(w, "TcpTracer status: %v\n", status)
}

func TracerReady(w http.ResponseWriter, req *http.Request) {
	status := "Not Ready"
	if tcptracer.IsReady() {
		status = "Ready"
	}
	fmt.Fprintf(w, "TcpTracer status: %v\n", status)
}

func main() {
	http.HandleFunc("/start", TracerStart)
	http.HandleFunc("/stop", TracerStop)
	http.HandleFunc("/status/stopped", TracerStopped)
	http.HandleFunc("/status/ready", TracerReady)

	fmt.Println("Running on http://127.0.0.1:8090/ (Press CTRL+C to quit)")
	http.ListenAndServe(":8090", nil)

	fmt.Printf("Bye\n")
}