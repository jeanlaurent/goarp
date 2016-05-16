package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type HostList struct {
	Hosts []Host
	lock  sync.Locker
}

func (h *HostList) Add(host Host) bool {
	h.lock.Lock()
	defer h.lock.Unlock()
	found := false
	for _, targetHost := range h.Hosts {
		if targetHost.Name == host.Name && targetHost.IP == host.IP {
			found = true
		}
	}
	if found {
		return false
	}
	h.Hosts = append(h.Hosts, host)
	return true
}

func (h *HostList) Print() {
	fmt.Println("Host found", len(h.Hosts))
	for _, host := range h.Hosts {
		fmt.Println(host)
	}
}

func NewHostList() *HostList {
	return &HostList{[]Host{}, &sync.Mutex{}}
}

type Host struct {
	Name string
	IP   string
	Mac  string
}

func (h Host) String() string {
	mac := ""
	if h.Mac != "" {
		mac = "[" + h.Mac + "]" + mac
	}
	return h.Name + "(" + h.IP + ")" + mac
}

func readHost(words []string) Host {
	host := Host{Name: words[0], IP: words[1][1 : len(words[1])-1], Mac: words[3]}
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
	var count int
	err := runArp(func(host Host) {
		count++
		if list.Add(host) {
			fmt.Println("Adding new host ", host)
		}
	})
	fmt.Println("Host found on this run", count)
	if err != nil {
		fmt.Println("Error on arp subsystem: ", err)
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
	for runNumber < 4 {
		time.Sleep(1 * time.Second)
	}
	close(quit)
	fmt.Println("=============")
	megaList.Print()
}
