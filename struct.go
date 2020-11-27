package godnsbl

import "sync"

// The main running service
type LookupService struct {
	DnsblListing []Dnsbl `yaml:"DroneListing" json:"DroneListing"`
	StartTime    int64
	TotalChecked uint64
}

// Each Dnsbl we add
type Dnsbl struct {
	Name       string            `yaml:"Name" json:"Name"`
	Address    string            `yaml:"Address" json:"Address"`
	Reply      map[string]string `yaml:"Reply" json:"Reply"`
	BanList    []int             `yaml:"BanList" json:"BanList"`
	BanMessage string            `yaml:"ManMessage" json:"BanMessage"`
}

// Standard data format for returning info from a Dnsbl
type DnsblReturn struct {
	IP      string
	Type    string
	Dnsbl   string
	Total   int
	Clear   int
	Message string
}

type dnsblCounter struct {
	ClearCount   int // Number that have not matched
	sync.RWMutex     // Mutex locker
}
