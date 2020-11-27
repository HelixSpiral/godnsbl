// Package godnsbl implements dnsbl lookups
package godnsbl

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// GetFirstDnsblReply gets the first reply from our lookup lists
func (l *LookupService) GetFirstDnsblReply(stringIP string) DnsblReturn {
	l.TotalChecked += 1
	reversedIP := reverseIP(stringIP)
	returnChan := make(chan DnsblReturn)
	counter := dnsblCounter{ClearCount: 0}

	for key := range l.DnsblListing {
		go func(key int) {
			lookup := fmt.Sprintf("%s%s", reversedIP, l.DnsblListing[key].Address) // Do the lookup for this BL
			lookupReply, err := net.LookupHost(lookup)                             // Grab the replies
			if err != nil && !strings.Contains(err.Error(), "no such host") {
				log.Fatal("Couldn't lookup the host:", err)
			}
			if len(lookupReply) == 0 || !replyMatch(lookupReply[0], l.DnsblListing[key].BanList) {
				counter.Lock()
				counter.ClearCount++
				counter.Unlock()
				return // We don't have any replies for the given IP/BL
			}

			returnChan <- DnsblReturn{
				IP:      stringIP,
				Type:    "BAN",
				Dnsbl:   l.DnsblListing[key].Name,
				Total:   len(l.DnsblListing),
				Message: strings.Replace(l.DnsblListing[key].BanMessage, "%IPADDR", stringIP, -1),
			}
		}(key)
	}

	// Loop until one of the following happen:
	// 1. We find a match in any of the dronebl services we're using
	// 2. We time out after 30 seconds
	// 3. We successfully pass all lookups
	for {
		select {
		case ok := <-returnChan:
			counter.RLock()
			ok.Clear = counter.ClearCount
			counter.RUnlock()
			return ok
		case <-time.After(l.timeout):
			counter.RLock()
			endcount := counter.ClearCount
			counter.RUnlock()
			return DnsblReturn{
				IP:    stringIP,
				Type:  "TIMEOUT",
				Dnsbl: "N/A",
				Total: len(l.DnsblListing),
				Clear: endcount,
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
					Clear: endcount,
				}
			}
		}
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
