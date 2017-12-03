package netstat

import (
	"github.com/ringtail/GOnetstat"
	"testing"
)

func TestNetstatLib(t *testing.T) {
	d := GOnetstat.Tcp()
	for index, p := range d {
		t.Logf("%d: %s %s\n", index, p.Name, p.State)
	}

	ud := GOnetstat.Udp()
	for index, p := range ud {
		t.Logf("%d: %s %s\n", index, p.Name, p.State)
	}
}
