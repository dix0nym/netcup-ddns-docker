package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	netcup "github.com/aellwein/netcup-dns-api/pkg/v1"
)

type IP struct {
	Query string
}

func getIP() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}

func update(customerNumber int, apiKey string, apiPassword string, domain string, subdomains map[string]bool) {
	var (
		dnsRecords      []netcup.DnsRecord
		existingRecords *[]netcup.DnsRecord
	)

	client := netcup.NewNetcupDnsClient(customerNumber, apiKey, apiPassword)

	session, err := client.Login()
	if err != nil {
		panic(err)
	}
	defer session.Logout()

	publicIP := getIP()

	if existingRecords, err = session.InfoDnsRecords(domain); err != nil {
		log.Printf("failed to get existing dns records: %v\n", err)
		return
	}

	for _, dnsRecord := range *existingRecords {
		if dnsRecord.Type != "A" {
			continue
		}
		if _, ok := subdomains[dnsRecord.Hostname]; ok {
			subdomains[dnsRecord.Hostname] = true
			if dnsRecord.Destination != publicIP {
				dnsRecord.Destination = publicIP
				dnsRecords = append(dnsRecords, dnsRecord)
				log.Printf("Set destination of A record %s to %s\n", dnsRecord.Hostname, dnsRecord.Destination)
			}
		}
	}

	for subdomain, exists := range subdomains {
		if !exists {
			dnsRecord := netcup.DnsRecord{Hostname: subdomain, Type: "A", Destination: publicIP}
			dnsRecords = append(dnsRecords, dnsRecord)
			log.Printf("Create new A record for %s to %s\n", dnsRecord.Hostname, dnsRecord.Destination)
		}
	}

	if len(dnsRecords) > 0 {
		if _, err := session.UpdateDnsRecords(domain, &dnsRecords); err != nil {
			log.Printf("failed to update dns records: %v\n", err)
			return
		}
	} else {
		log.Println("No updates found")
	}
}

func main() {
	var (
		domain         string
		subdomains     []string
		interval       int
		customerNumber int
		apiKey         string
		apiPassword    string
		err            error
	)

	if domain = os.Getenv("DOMAIN"); domain == "" {
		panic("env DOMAIN empty")
	}

	if subdomains = strings.Split(os.Getenv("SUBDOMAINS"), ","); len(subdomains) == 0 {
		panic("env SUBDOMAINS empty")
	}

	if interval, err = strconv.Atoi(os.Getenv("INTERVAL")); err != nil {
		log.Panicf("env INTERVAL empty or not an int: %v", err)
	}

	if customerNumber, err = strconv.Atoi(os.Getenv("NETCUP_CUSTOMER_NUMBER")); err != nil {
		log.Panicf("env NETCUP_CUSTOMER_NUMBER empty or not an int: %v", err)
	}

	if apiKey = os.Getenv("NETCUP_API_KEY"); apiKey == "" {
		panic("env NETCUP_API_KEY empty")
	}

	if apiPassword = os.Getenv("NETCUP_API_PASSWORD"); apiPassword == "" {
		panic("env NETCUP_API_PASSWORD empty")
	}

	subdomainMap := make(map[string]bool)
	for _, subdomain := range subdomains {
		subdomainMap[subdomain] = false
	}

	if interval != 0 {
		for {
			update(customerNumber, apiKey, apiPassword, domain, subdomainMap)
			time.Sleep(time.Duration(interval) * time.Second)
		}
	} else {
		update(customerNumber, apiKey, apiPassword, domain, subdomainMap)
	}
}
