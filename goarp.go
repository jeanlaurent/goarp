package main

// Find IP Range
// ifconfig en0 | grep -v inet6 | grep inet | cut -d ' ' -f 2
// 10.0.1.25

// Discard name with ?

// find more info on a host
//nmap -sn -oG - -e en0 10.0.1.1/24
//# Nmap 7.01 scan initiated Sat May 28 19:00:23 2016 as: nmap -sn -oG - -e en0 10.0.1.1/24
//Host: 10.0.1.1 (livebox.home)	Status: Up
//Host: 10.0.1.2 (timecapsule.home)	Status: Up
//Host: 10.0.1.18 (ll.home)	Status: Up
//Host: 10.0.1.20 (t.home)	Status: Up
//Host: 10.0.1.25 (ll.home)	Status: Up
//Host: 10.0.1.69 (ipad-morlhon.home)	Status: Up
//# Nmap done at Sat May 28 19:00:31 2016 -- 256 IP addresses (6 hosts up) scanned in 8.24 seconds

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// HostList maintain a list of Hosts
type HostList struct {
	Hosts []*Host
	lock  sync.Locker
}

// Add a Host to the list if not already there.
func (h *HostList) Add(host Host) bool {
	h.lock.Lock()
	defer h.lock.Unlock()
	if host.Name == "?" {
		return false
	}
	found := false
	for _, targetHost := range h.Hosts {
		if targetHost.Name == host.Name && targetHost.IP == host.IP {
			targetHost.LastSeen = time.Now() // wont work as this is a copy of the object
			found = true
		}
	}
	if found {
		return false
	}
	h.Hosts = append(h.Hosts, &host)
	return true
}

// Print to stdout the hostlist
func (h *HostList) Print() {
	fmt.Println("Host found", len(h.Hosts))
	for _, host := range h.Hosts {
		fmt.Println(host)
	}
}

// Size returns the number of element of the list
func (h *HostList) Size() int {
	return len(h.Hosts)
}

// NewHostList created an empty HostList
func NewHostList() *HostList {
	return &HostList{[]*Host{}, &sync.Mutex{}}
}

// Host is a registered host
type Host struct {
	Name      string
	IP        string
	Mac       string
	FirstSeen time.Time
	LastSeen  time.Time
}

func (h Host) String() string {
	mac := ""
	if h.Mac != "" {
		mac = "[" + h.Mac + "]" + mac
	}
	lastSeen := ""
	if !h.LastSeen.IsZero() {
		lastSeen = " ==> " + h.LastSeen.Format(time.UnixDate)
	}
	return h.Name + "(" + h.IP + ")" + mac + " @ " + h.FirstSeen.Format(time.UnixDate) + lastSeen
}

func readHost(words []string) Host {
	host := Host{Name: words[0], IP: words[1][1 : len(words[1])-1], Mac: words[3], FirstSeen: time.Now()}
	if host.Mac[0] == '(' {
		host.Mac = ""
	}
	return host
}

func runArp(onHost func(Host)) error {
	cmdName := "arp"
	cmdArgs := []string{"-i", "en0", "-a"}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		return err
	}

	lines := strings.Split(string(cmdOut), "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		if len(words) > 3 {
			onHost(readHost(words))
		}
	}
	return nil
}

func run(list *HostList) {
	err := runArp(func(host Host) {
		if list.Add(host) {
			fmt.Println("Adding new host", host)
		}
	})
	fmt.Println("Host detect on this run", list.Size())
	if err != nil {
		fmt.Println("Error on arp subsystem:", err)
	}
}

func main() {
	megaList := NewHostList()
	run(megaList)

	var runNumber int
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				runNumber++
				run(megaList)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	for runNumber < 20 { // 4 is 3
		time.Sleep(1 * time.Second)
	}
	close(quit)
	fmt.Println("=============")
	megaList.Print()
}
