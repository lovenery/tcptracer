package tracer

import(
	"net"
	"testing"
	"github.com/lovenery/tcptracer/pkg/tcptracer"
)

func TestIpIntToByte(t *testing.T) {
	want := net.ParseIP("10.0.2.15")
	have := tcptracer.IpIntToByte(251789322)

	t.Log(have.String())
	t.Log(want.String())

	if have.String() != want.String() {
		t.Errorf("Start(): have: '%v', want: '%v'", have, want)
	}
}

func TestAll(t *testing.T) {
	for {
		if tcptracer.IsStopped() {
			break
		}
	}

	tcptracer.Start()

	for {
		if tcptracer.IsReady() {
			break
		}
	}

	tcptracer.Stop()

	for {
		if tcptracer.IsStopped() {
			break
		}
	}
}