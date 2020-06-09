// This file is part of myhttp

// Licensed under the BSD-3-Clause license:
// https://opensource.org/licenses/BSD-3-Clause
// Copyright (c) 2020, Pablo S. Blum de Aguiar

package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
)

func parseArguments(maxParallel uint64) (uint64, []string, error) {
	parallel := flag.Uint64("parallel", 10, "limit the number of parallel requests")
	flag.Usage = usage
	flag.Parse()
	if (*parallel) < 1 {
		return 0, nil, fmt.Errorf(`invalid value "%d" for flag -parallel`, (*parallel))
	}
	if (*parallel) > maxParallel {
		log.Printf(
			`Warning: value "%d" for flag -parallel is too big, using default max %d`,
			(*parallel),
			maxParallel,
		)
		parallel = &maxParallel
	}
	urls, err := parseURLs()
	return *parallel, urls, err
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: myhttp [flags] [url ...]\n")
	flag.PrintDefaults()
}

// parseURLs parses URLs provided as arguments to the program
func parseURLs() ([]string, error) {
	if flag.NArg() == 0 {
		return nil, fmt.Errorf("no url provided")
	}
	urls := make([]string, 0, flag.NArg())
	for i := 0; i < flag.NArg(); i++ {
		if url, err := parseURL(flag.Arg(i)); err == nil {
			urls = append(urls, url)
		} else {
			log.Printf("Warning: %s is invalid (%s)\n", flag.Arg(i), err)
		}
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("no valid url provided")
	}
	return urls, nil
}

// parseURLs parses a URL ensuring it's valid and prepending a protocol if needed
func parseURL(myURL string) (string, error) {
	myURL = strings.TrimSpace(myURL)
	if strings.HasPrefix(myURL, "//") {
		myURL = "http:" + myURL
	}
	parsed, err := url.Parse(myURL)
	if err != nil {
		return myURL, err
	}
	if parsed.Scheme == "" {
		parsed, err = url.Parse("http://" + myURL)
		if err != nil {
			return myURL, err
		}
	}
	myURL = parsed.String()
	return strings.TrimSuffix(myURL, "/"), nil
}
