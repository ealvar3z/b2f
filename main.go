// b2f.go
// bitwarden2factotum is a utility to convert a Bitwarden CSV export
// into a factotum ctl commands.
// Usage:
//   go run b2f.go -input vault.csv > facts.ctl
//   cat facts.ctl > /mnt/factotum/ctl
// Or use the -apply flag to write directly:
//   go run b2f.go -input vault.csv -apply

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// flags
	input := flag.String("input", "", "Path to Bitwarden CSV export")
	mount := flag.String("mount", "/mnt/factotum", "Factotum mountpoint")
	apply := flag.Bool("apply", false, "Apply directly to factotum ctl file")
	outFile := flag.String("out", "", "Output file (default stdout)")
	flag.Parse()

	if *input == "" {
		log.Fatal("Provide -input path to Bitwarden CSV export")
	}

	f, err := os.Open(*input)
	if err != nil {
		log.Fatalf("Error operning input file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	head, err := r.Read()
	if err != nil {
		log.Fatalf("Error reading header: %v", err)
	}

	idxURI := indexOf(head, "login_uri")
	idxUser := indexOf(head, "login_username")
	idxPass := indexOf(head, "login_password")
	if idxURI < 0 || idxUser < 0 || idxPass < 0 {
		log.Fatalf("Missing expected columns: login_uri, login_username, login_password")
	}

	var out io.Writer = os.Stdout
	if *apply {
		ctl := filepath.Join(*mount, "ctl")
		file, err := os.OpenFile(ctl, os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("Error opening %s: %v", ctl, err)
		}
		defer file.Close()
		out = file
	} else if *outFile != "" {
		file, err := os.Create(*outFile)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer file.Close()
		out = file
	}

	valid := regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

	// process records
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading record: %v", err)
		}

		uri := rec[idxURI]
		user := rec[idxUser]
		pass := rec[idxPass]
		if uri == "" || user == "" || pass == "" {
			continue
		}

		host := parseHost(uri)
		if host == "" {
			continue
		}

		fmt.Fprintf(out, "key proto=pass service=%s user=%s", quoteVal(host, valid), quoteVal(user, valid))
		fmt.Fprintf(out, "!password=%s", singleQuote(pass))
	}
}

func indexOf(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}
	return -1
}

func parseHost(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return ""
	}
	return u.Hostname()
}

func singleQuote(s string) string {
	esc := strings.ReplaceAll(s, "'", "''")
	return "'" + esc + "'"
}

func quoteVal(val string, valid *regexp.Regexp) string {
	if valid.MatchString(val) {
		return val
	}
	return singleQuote(val)
}
