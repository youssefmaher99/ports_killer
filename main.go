package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func kill(startRange int, endRange int) error {
	cmd := exec.Command("/usr/bin/lsof", "-i")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running lsof: %v", err)
	}
	processIDs := make(map[string]struct{})
	lines := strings.Split(string(data), "\n")
	if len(lines) <= 1 {
		return fmt.Errorf("no open ports in this range")
	}

	for i := 1; i < len(lines)-1; i++ {
		pid := strings.Fields(lines[i])[1]
		port := strings.Split(strings.Fields(lines[i])[8], ":")[1]
		port_int, err := strconv.Atoi(port)
		if err != nil {
			continue
		}

		if port_int >= startRange && port_int <= endRange {
			processIDs[pid] = struct{}{}
		}
	}

	for pid := range processIDs {
		cmd := exec.Command("/usr/bin/kill", "-9", pid)
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Killed")
	}

	return nil
}

func main() {
	start := flag.Int("start", -1, "start port range")
	end := flag.Int("end", -1, "end port range")
	flag.Parse()
	if *start == -1 || *end == -1 {
		panic("Invalid arguments --start or --end must in range [1-65535] ")
	}
	err := kill(*start, *end)
	if err != nil {
		log.Fatal(err)
	}
}
