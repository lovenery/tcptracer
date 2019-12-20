package tcptracer

import(
	"net"
	"testing"
)

func TestIpIntToByte(t *testing.T) {
	want := net.ParseIP("10.0.2.15")
	have := IpIntToByte(251789322)

	t.Log(have.String())
	t.Log(want.String())

	if have.String() != want.String() {
		t.Errorf("Start(): have: '%v', want: '%v'", have, want)
	}
}