package tcptracer

import (
	"C"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
	bpf "github.com/iovisor/gobpf/bcc"
)

func IpIntToByte(ip uint32) net.IP {
	result := make(net.IP, 4)
			result[0] = byte(ip)
			result[1] = byte(ip >> 8)
			result[2] = byte(ip >> 16)
			result[3] = byte(ip >> 24)
	return result
}

func TcpTypeIntToString(TcpType uint32) string {
	result := "unknown"
	if TcpType == 1 {
		result = "connect"
	} else if TcpType == 2 {
		result = "accept"
	} else if TcpType == 3 {
		result = "close"
	}
	return result
}

type tcpIpv4Event struct {
	TSns        uint64 // Current TimeStamp in nanoseconds
	TcpType     uint32
	Pid         uint32
	Comm        [16]byte // TASK_COMM_LEN=16
	IpVer       uint8
	Padding     [3]byte
	Saddr       uint32
	Daddr       uint32
	Sport       uint16
	Dport       uint16
	Netns       uint32
}

var IsTracerDoneSig = make(chan bool, 1)
var IsTracerStopped bool = true
var IsTracerReady bool = false

func receiveChan(channel chan []byte) {
	var event tcpIpv4Event
	fmt.Println("[lib] Tcp Tracer is Ready")
	IsTracerReady = true
	for {
		data := <-channel
		err := binary.Read(bytes.NewBuffer(data), bpf.GetHostByteOrder(), &event)
		if err != nil {
			fmt.Printf("failed to decode received data: %s\n", err)
			continue
		}

		fmt.Printf("-------------------\n")
		log.Print()
		fmt.Printf("TSns   : %d \n", event.TSns)
		fmt.Printf("Type   : %s \n", TcpTypeIntToString(event.TcpType))
		fmt.Printf("PID    : %d \n", event.Pid)
		Comm := (*C.char)(unsafe.Pointer(&event.Comm))
		fmt.Printf("COMM   : %s \n", C.GoString(Comm))
		fmt.Printf("IP     : IPv%d \n", event.IpVer)
		fmt.Printf("SADDR  : %s \n", IpIntToByte(event.Saddr))
		fmt.Printf("SADDR  : %v \n", event.Saddr)
		fmt.Printf("DADDR  : %s \n", IpIntToByte(event.Daddr))
		fmt.Printf("SPORT  : %d \n", event.Sport)
		fmt.Printf("DPORT  : %d \n", event.Dport)
		fmt.Printf("NETNS  : %d \n", event.Netns)
	}
}

func StartMain() {
	// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
	_, b, _, _ := runtime.Caller(0)
	basePath   := filepath.Dir(b)
	newPath := filepath.Join(basePath, "tcptracer.bt")

	sourceByte, err := ioutil.ReadFile(newPath)
	if err != nil {
		log.Fatal(err)
	}
	source := string(sourceByte)
	m := bpf.NewModule(source, []string{})
	defer m.Close()

	kprobe, err := m.LoadKprobe("trace_connect_v4_entry")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load trace_connect_v4_entry: %s\n", err)
		os.Exit(1)
	}
	m.AttachKprobe("tcp_v4_connect", kprobe, 0)
	kprobe, err = m.LoadKprobe("trace_connect_v4_return")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load trace_connect_v4_return: %s\n", err)
		os.Exit(1)
	}
	m.AttachKretprobe("tcp_v4_connect", kprobe, 0)
	kprobe, err = m.LoadKprobe("trace_tcp_set_state_entry")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load trace_tcp_set_state_entry: %s\n", err)
		os.Exit(1)
	}
	m.AttachKprobe("tcp_set_state", kprobe, 0)
	kprobe, err = m.LoadKprobe("trace_close_entry")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load trace_close_entry: %s\n", err)
		os.Exit(1)
	}
	m.AttachKprobe("tcp_close", kprobe, 0)
	kprobe, err = m.LoadKprobe("trace_accept_return")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load trace_close_entry: %s\n", err)
		os.Exit(1)
	}
	m.AttachKretprobe("inet_csk_accept", kprobe, 0)

	table := bpf.NewTable(m.TableId("tcp_ipv4_event"), m)
	channel := make(chan []byte)
	perfMap, err := bpf.InitPerfMap(table, channel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init perf map: %s\n", err)
		os.Exit(1)
	}

	go receiveChan(channel)

	perfMap.Start()
	<-IsTracerDoneSig
	perfMap.Stop()
}

func Start() error {
	if IsTracerStopped == true {
		IsTracerDoneSig = make(chan bool, 1)
		IsTracerStopped = false
		IsTracerReady = false
	} else {
		return errors.New("Tcp Tracer is running")
	}

	go StartMain()

	return nil
}

func Stop() {
	IsTracerDoneSig <- true
	fmt.Println("[lib] Tcp Tracer is Stopped")
	IsTracerStopped = true
	IsTracerReady = false
}

func IsStopped() bool {
	return IsTracerStopped
}

func IsReady() bool {
	return IsTracerReady
}