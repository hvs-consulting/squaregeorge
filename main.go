package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	printDisclaimer()
	var outputStdout bool
	var f *os.File
	domainMap := make(map[string][]*net.MX)

	if len(os.Args) > 1 {
		var err error
		// Read from file
		f, err = os.Open(os.Args[1])
		if err != nil {
			panic(fmt.Sprintf("Error opening file: %s", os.Args[1]))
		}
		outputStdout = false
	} else {
		// Read from stdin
		f = os.Stdin
		outputStdout = true
	}

	var mx []*net.MX
	var ok bool
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		target := scanner.Text()
		// Filter input
		sanitized, err := sanitizeInput(target)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: Error occured: %s", target, err.Error()))
		}
		// Check if the domain has already been looked up
		mx, ok = domainMap[sanitized]
		if !ok {
			// if not, do the lookup
			mx, err = net.LookupMX(sanitized)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: Error occured: %s", target, err.Error()))
			}
			domainMap[sanitized] = mx
		}
	}

	processOutput(outputStdout, domainMap)
	if !outputStdout {
		fmt.Println("Done. Waiting 3 seconds before closing the window.")
		time.Sleep(3 * time.Second)
	}
}

func printDisclaimer() {
	fmt.Fprintln(os.Stderr, "squaregeorge - Identify mailservers in charge for delivery to a given set of addresses")
	fmt.Fprintln(os.Stderr, "By Michael Eder, HvS-Consulting AG")
	fmt.Fprintln(os.Stderr, "https://www.hvs-consulting.de/  https://twitter.com/michael_eder_")
	fmt.Fprintln(os.Stderr, "")
}

// sanitizeInput does some rudimentary input checking
func sanitizeInput(s string) (string, error) {
	// Remove trailing spaces
	trimmed := strings.TrimSpace(s)

	// Check if the line contains a '@'. If yes, strip everything up to it
	if i := strings.LastIndex(trimmed, "@"); i != -1 {
		// avoid the case where '@' is the last character
		if i < len(trimmed) {
			trimmed = trimmed[i+1:]
		} else {
			return "", fmt.Errorf("Invalid name: %s", s)
		}
	}

	// Each domain should at least contain a dot, so do this as final check
	if strings.Contains(trimmed, ".") {
		return trimmed, nil
	}
	return "", fmt.Errorf("Invalid name: %s", s)

}

func processOutput(writeStdout bool, r map[string][]*net.MX) {
	var w *csv.Writer
	if writeStdout {
		w = csv.NewWriter(os.Stdout)
	} else {
		f, err := os.Create("resolved_mx_domains.csv")
		if err != nil {
			panic(fmt.Errorf("Opening the file resolved_my_domains.csv failed: %s", err.Error()))
		}
		defer f.Close()
		b := bufio.NewWriter(f)
		w = csv.NewWriter(b)
	}

	// Goal is to make this easy with Excel on Windows
	w.Comma = ';'
	w.UseCRLF = true

	// Write the data to the CSV file
	w.Write([]string{"Domain", "MX", "Preference"})
	for domain, mxes := range r {
		for _, mx := range mxes {
			err := w.Write([]string{domain, mx.Host, strconv.FormatUint((uint64)(mx.Pref), 10)})
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to write CSV: %s", err.Error()))
			}
		}
	}
	w.Flush()
}
