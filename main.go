// This file is part of myhttp

// Licensed under the BSD-3-Clause license:
// https://opensource.org/licenses/BSD-3-Clause
// Copyright (c) 2020, Pablo S. Blum de Aguiar

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	parallel, urls, err := parseArguments(getMaxParallel())
	if err != nil {
		log.Printf("Error: %s\n", err)
		usage()
		os.Exit(1)
	}
	client := newClientOfMine()
	pageFeed := client.request(urls, parallel)
	for pageResult := range pageFeed {
		fmt.Printf("%s\n", pageResult)
	}
}
