package godnsbl

import (
	"sync"
	"time"
)

// LookupService is the main running service
type LookupService struct {
	DnsblTimeout string  `yaml:"DnsblTimeout" json:"DnsblTimeout"`
	DnsblListing []Dnsbl `yaml:"DroneListing" json:"DroneListing"`
	StartTime    int64
	TotalChecked uint64

	// We parse DnsblTimeout and store it in an unexported time.Duration so we only have to parse it once
	timeout time.Duration
}

// Dnsbl stores the data for each Dnbsl we have
type Dnsbl struct {
	Name         string            `yaml:"Name" json:"Name"`
	Address      string            `yaml:"Address" json:"Address"`
	Reply        map[string]string `yaml:"Reply" json:"Reply"`
	AllowList    []int             `yaml:"AllowList" json:"AllowList"`
	BlockMessage string            `yaml:"BlockMessage" json:"BlockMessage"`
}

// DnsblReturn is the standard data format for returning info from a Dnsbl
type DnsblReturn struct {
	IP      string
	Type    string
	Dnsbl   string
	Message string
}

type dnsblCounter struct {
	ClearCount   int // Number that have not matched
	sync.RWMutex     // Mutex locker
}
