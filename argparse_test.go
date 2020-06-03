// This file is part of myhttp

// Licensed under the BSD-3-Clause license:
// https://opensource.org/licenses/BSD-3-Clause
// Copyright (c) 2020, Pablo S. Blum de Aguiar

package main

import (
	"flag"
	"testing"
)

func TestParseURL(t *testing.T) {
	testCases := []struct {
		url      string
		expected string
	}{
		{
			"golang.org",
			"http://golang.org",
		}, {
			"http://duckduckgo.com/",
			"http://duckduckgo.com",
		}, {
			"http://www.globo.com",
			"http://www.globo.com",
		}, {
			"https://g1.globo.com",
			"https://g1.globo.com",
		}, {
			"//www.debian.org",
			"http://www.debian.org",
		}, {
			" foo.html",
			"http://foo.html",
		},
	}
	for _, testCase := range testCases {
		parsed, err := parseURL(testCase.url)
		if err != nil {
			t.Error(err)
		}
		if parsed != testCase.expected {
			t.Errorf("%s should be equal to %s", parsed, testCase.expected)
		}
	}
}

func TestParseURLFails(t *testing.T) {
	URLs := []string{"http://[fe80::1%en0]", "googl%65.com"}
	for _, url := range URLs {
		_, err := parseURL(url)
		if err == nil {
			t.Errorf("Should fail with %s", url)
		}
	}
}

func TestParseURLs(t *testing.T) {
	testCases := []struct {
		URLs        []string
		expectedLen int
	}{
		{
			[]string{
				"golang.org",
				"http://duckduckgo.com/",
				"http://www.globo.com",
				"https://g1.globo.com",
				"//www.debian.org",
				" foo.html",
			},
			6,
		}, {
			[]string{"golang.org", "http://[fe80::1%en0]", "googl%65.com"},
			1,
		},
	}
	for _, testCase := range testCases {
		flag.CommandLine.Parse(testCase.URLs)
		parsed, err := parseURLs()
		if err != nil {
			t.Error(err)
		}
		if len(parsed) != testCase.expectedLen {
			t.Errorf("%q should have %d item(s)s", parsed, testCase.expectedLen)
		}
	}
}

func TestParseURLsFails(t *testing.T) {
	testCases := [][]string{{}, {"http://[fe80::1%en0]", "googl%65.com"}}
	for _, URLs := range testCases {
		flag.CommandLine.Parse(URLs)
		_, err := parseURLs()
		if err == nil {
			t.Errorf("Should fail with %q", URLs)
		}
	}
}
