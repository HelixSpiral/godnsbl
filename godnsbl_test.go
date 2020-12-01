package godnsbl_test

import (
	"testing"

	"github.com/HelixSpiral/godnsbl"
)

func TestLookupIP(t *testing.T) {
	dnsbl := godnsbl.NewLookupService()

	dnsbl.DnsblListing = append(dnsbl.DnsblListing, godnsbl.Dnsbl{
		Name:         "Test - DroneBL",
		Address:      ".dnsbl.dronebl.org",
		BlockList:    []int{1},
		BlockMessage: "%IPADDR found",
	})

	reply := dnsbl.LookupIP("127.0.0.2")

	if reply.Type != "BLOCK" {
		t.Errorf("Expended a ban, got: %+v\r\n", reply)
	}

	reply = dnsbl.LookupIP("2001:4860:4860::8888")

	if reply.Type != "CLEAR" {
		t.Errorf("Expended a clear, got: %+v\r\n", reply)
	}
}
