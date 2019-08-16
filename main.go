package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	cidrman "github.com/EvilSuperstars/go-cidrman"
)

func main() {

	hostPrefix := flag.Int("hostPrefix", 32, "prefix length for host routes summary")
	listURL := flag.String("listURL", "https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv", "URL to get list of banned resources")
	flag.Parse()

	resp, err := http.Get(*listURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	dump, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	split := bytes.Split(dump, []byte("\n"))

	var prefixes []string
	for _, line := range split {
		IPs := bytes.Split(line, []byte(";"))
		IPs = bytes.Split(IPs[0], []byte(" | "))
		for _, ip := range IPs {
			prefixes = append(prefixes, string(ip))
		}
	}

	var sanitized []string
	for _, prefix := range prefixes {
		if strings.Contains(prefix, "/") {
			sanitized = append(sanitized, prefix)
			continue
		}
		if net.ParseIP(prefix).To4() != nil {
			sanitized = append(sanitized, fmt.Sprintf("%s/%d", prefix, *hostPrefix))
			continue
		}
	}

	merged, err := cidrman.MergeCIDRs(sanitized)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range merged {
		fmt.Println("route " + v + " reject;")
	}
	log.Println(len(sanitized), "=>>", len(merged))
}
