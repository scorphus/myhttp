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
	_, urls, err := parseArguments()
	if err != nil {
		log.Printf("Error: %s\n", err)
		usage()
		os.Exit(1)
	}
	for _, url := range urls {
		fmt.Printf("%s\n", url)
	}
}
