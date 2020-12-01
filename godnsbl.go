// Package godnsbl implements dnsbl lookups
package godnsbl

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// LookupIP gets the first reply from our lookup lists
func (l *LookupService) LookupIP(stringIP string) DnsblReturn {
	l.TotalChecked++
	reversedIP := reverseIP(stringIP)
	returnChan := make(chan DnsblReturn)
	counter := dnsblCounter{ClearCount: 0}

	lookupCtx, lookupCancel := context.WithTimeout(context.Background(), l.timeout)
	defer lookupCancel()

	for key := range l.DnsblListing {
		go l.dnsblLookup(lookupCtx, returnChan, stringIP, reversedIP, key)
	}

	// Loop until one of the following happen:
	// 1. We find a match in any of the dronebl services we're using
	// 2. We time out after the specified timeout
	// 3. We successfully pass all lookups
	for {
		select {
		case ok := <-returnChan:
			counter.Lock()
			counter.ClearCount++
			counter.Unlock()
			if ok.Type == "BLOCK" {
				return ok
			}
		case <-lookupCtx.Done():
			return DnsblReturn{
				IP:    stringIP,
				Type:  "TIMEOUT",
				Dnsbl: "N/A",
			}
		default:
			counter.RLock()
			endcount := counter.ClearCount
			counter.RUnlock()
			if endcount == len(l.DnsblListing) {
				return DnsblReturn{
					IP:    stringIP,
					Type:  "CLEAR",
					Dnsbl: "N/A",
				}
			}
		}
	}
}

// LookupIPGetAll gets replies from all Dnsbls, not just the first to reply with a block
func (l *LookupService) LookupIPGetAll(stringIP string) []DnsblReturn {
	l.TotalChecked++
	reversedIP := reverseIP(stringIP)
	returnChan := make(chan DnsblReturn)
	var replyList []DnsblReturn

	lookupCtx, lookupCancel := context.WithTimeout(context.Background(), l.timeout)
	defer lookupCancel()

	for key := range l.DnsblListing {
		go l.dnsblLookup(lookupCtx, returnChan, stringIP, reversedIP, key)
	}

	// Wait until we get all the replies or we time out
	for {
		select {
		case <-lookupCtx.Done(): // We timed out, return what we have
			return replyList
		case reply := <-returnChan:
			replyList = append(replyList, reply)
			if len(replyList) == len(l.DnsblListing) {
				return replyList
			}
		}
	}
}

func (l *LookupService) dnsblLookup(lookupCtx context.Context, returnChan chan<- DnsblReturn, stringIP, reversedIP string, key int) {
	returnDnsbl := DnsblReturn{
		IP:      stringIP,
		Type:    "CLEAR",
		Dnsbl:   l.DnsblListing[key].Name,
		Message: "IP not found",
	}
	lookup := fmt.Sprintf("%s%s", reversedIP, l.DnsblListing[key].Address)
	lookupReply, err := net.LookupHost(lookup)
	if err != nil && !err.(*net.DNSError).IsNotFound {
		log.Fatal("Couldn't lookup the host:", err)
	}
	if len(lookupReply) == 0 || !replyMatch(lookupReply[0], l.DnsblListing[key].BlockList) {
		returnChan <- returnDnsbl
		return
	}

	// Replace %IPADDR with the IP we're checking
	returnMessage := strings.Replace(l.DnsblListing[key].BlockMessage, "%IPADDR", stringIP, -1)

	// Lookup if we have a known index for the given reply, if we do tack it on the end of the message.
	replyIndex := strings.LastIndex(lookupReply[0], ".") + 1
	if l.DnsblListing[key].Reply[lookupReply[0][replyIndex:]] != "" {
		returnMessage = fmt.Sprintf("%s (%s)", returnMessage, l.DnsblListing[key].Reply[lookupReply[0][replyIndex:]])
	}

	returnDnsbl.Type = "BLOCK"
	returnDnsbl.Message = returnMessage

	select {
	case <-lookupCtx.Done():
		return
	case returnChan <- returnDnsbl:
		return
	}
}

// replyMatch checks to see if the string we got has a reply that's in our list of matches
// that we want to ban
func replyMatch(check string, list []int) bool {
	checkstrip := strings.LastIndex(check, ".") + 1   // Strip all but the last part
	checkint, err := strconv.Atoi(check[checkstrip:]) // Convert it to an int
	if err != nil {
		log.Fatal("Couldn't make an int: ", err)
	}

	// Loop for our ban list
	for _, value := range list {
		if value == checkint {
			return true // Ban it!
		}
	}

	return false // It's approved, let it pass.

}

// reverseIP takes the IP we're given and reverse it
func reverseIP(IP string) string {
	var stringSplitIP []string

	// Do IPvX-specific splitting
	if net.ParseIP(IP).To4() != nil {
		stringSplitIP = strings.Split(IP, ".") // Split into 4 groups
	} else {
		stringSplitIP = strings.Split(IP, ":") // Split into however many groups

		/* Due to IPv6 lookups being different than IPv4 we have an extra check here
		We have to expand the :: and do 0-padding if there are less than 4 digits */
		for key := range stringSplitIP {
			if len(stringSplitIP[key]) == 0 { // Found the ::
				stringSplitIP[key] = strings.Repeat("0000", 8-strings.Count(IP, ":"))
			} else if len(stringSplitIP[key]) < 4 { // 0-padding needed
				stringSplitIP[key] = strings.Repeat("0", 4-len(stringSplitIP[key])) + stringSplitIP[key]
			}
		}

		// We have to join what we have and split it again to get all the letters split individually
		stringSplitIP = strings.Split(strings.Join(stringSplitIP, ""), "")
	}

	// Reverse the IP, join by . and return it
	for x, y := 0, len(stringSplitIP)-1; x < y; x, y = x+1, y-1 {
		stringSplitIP[x], stringSplitIP[y] = stringSplitIP[y], stringSplitIP[x] // Reverse the groups
	}

	return strings.Join(stringSplitIP, ".") // Return the IP.
}
