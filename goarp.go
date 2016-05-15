package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

func main() {
	var list []Host
	cmdName := "arp"
	cmdArgs := []string{"-i", "en0", "-a"}

	cmdOut, err := exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		fmt.Println("Error on arp subsystem: ", err)
		os.Exit(1)
	}

	lines := strings.Split(string(cmdOut), "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		if len(words) > 3 {
			list = append(list, readHost(words))
		}
	}

	for _, host := range list {
		fmt.Println(host)
	}

}
