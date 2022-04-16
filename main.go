package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %s\n", err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	if mxRecords, err := net.LookupMX(domain); err != nil {
		log.Printf("Error: %v\n", err)
	} else if len(mxRecords) > 0 {
		hasMX = true
	}

	if txtRecords, err := net.LookupTXT(domain); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		for _, record := range txtRecords {
			if strings.HasPrefix(record, "v=spf1") {
				hasSPF = true
				spfRecord = record
				break
			}
		}
	}

	if dmarcRecords, err := net.LookupTXT("_dmarc." + domain); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		for _, record := range dmarcRecords {
			if strings.HasPrefix(record, "v=DMARC1") {
				hasDMARC = true
				dmarcRecord = record
				break
			}
		}
	}

	if !hasSPF {
		spfRecord = "nil"
	}
	if !hasDMARC {
		dmarcRecord = "nil"
	}
	fmt.Printf("%v, %v, %v, %v, %v, %v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
