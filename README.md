# README

This is a simple package for doing dronebl lookups that was originally part of an IRC bot written in ~2015. Finally got around to separating this logic out into it's own package.

You can create a new service with:
```go
blService := godnsbl.NewLookupServiceWithConfig("config.yml")
```
or with
```go
blService := godnsbl.NewLookupService()

blService.DnsblListing = append(blService.DnsblListing, godnsbl.Dnsbl{
    Name: "DroneBL",
    Address: ".dnsbl.dronebl.org",
    Reply: map[string]string{
        "2": "Test",
    },
    BlockList: []int{1},
    BlockMessage: "Example block msg, your IP is: %IPADDR",
})
```

You can then use the lookup service:
```go
lookup := blService.GetFirstDnsblReply("127.0.0.2")
```