package godnsbl_test

import (
	"testing"

	"github.com/HelixSpiral/godnsbl"
)

func TestGetFirstDnsblReply(t *testing.T) {
	dnsbl := godnsbl.NewLookupService()

	dnsbl.DnsblListing = append(dnsbl.DnsblListing, godnsbl.Dnsbl{
		Name:       "Test - DroneBL",
		Address:    ".dnsbl.dronebl.org",
		BanList:    []int{1},
		BanMessage: "%IPADDR found",
	})

	reply := dnsbl.GetFirstDnsblReply("127.0.0.2")

	if reply.Type != "BAN" {
		t.Errorf("Expended a ban, got: %+v\r\n", reply)
	}

	reply = dnsbl.GetFirstDnsblReply("2001:4860:4860::8888")

	if reply.Type != "CLEAR" {
		t.Errorf("Expended a clear, got: %+v\r\n", reply)
	}
}
