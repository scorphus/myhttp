// This file is part of myhttp

// Licensed under the BSD-3-Clause license:
// https://opensource.org/licenses/BSD-3-Clause
// Copyright (c) 2020, Pablo S. Blum de Aguiar

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDoGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("I love pancakes"))
	}))
	defer server.Close()
	client := newClientOfMineWithHTTPClient(server.Client())
	actual, err := client.doGet(server.URL)
	if err != nil {
		t.Errorf("Got an unexpected error: %s", err)
	}
	if string(actual) != "I love pancakes" {
		t.Errorf("Got an unexpected response: %s", actual)
	}
}

func TestDoGetSucceedsWithBadStatus(t *testing.T) {
	statuses := []int{
		http.StatusBadRequest,
		http.StatusInternalServerError,
		http.StatusNotFound,
		http.StatusUnauthorized,
		http.StatusRequestTimeout,
	}
	for _, status := range statuses {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(status)
		}))
		client := newClientOfMineWithHTTPClient(server.Client())
		_, err := client.doGet(server.URL)
		if err != nil {
			t.Errorf("Should not fail with status %d", status)
		}
		server.Close()
	}
}

func TestDoGetFailsWithBadURLs(t *testing.T) {
	urls := []string{"\n", "\t", "%&", ":", "%30:%31:%32:%33", "/foo.html"}
	client := newClientOfMineWithHTTPClient(&http.Client{})
	for _, url := range urls {
		_, err := client.doGet(url)
		if err == nil {
			t.Errorf("Should fail for %s", url)
		}
	}
}

type roundTripper func(req *http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func roundTripperOfSuccess(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("Response from " + req.URL.String())),
		Header:     make(http.Header),
	}, nil
}

func roundTripperOfFailure(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("some error")
}

func TestDoGetFailsWithError(t *testing.T) {
	client := newClientOfMineWithHTTPClient(&http.Client{
		Transport: roundTripper(roundTripperOfFailure),
	})
	_, err := client.doGet("http://coffeegeek.com")
	if err == nil || !strings.Contains(err.Error(), "some error") {
		t.Error("Should fail with some error")
	}
}

func TestDoGetPage(t *testing.T) {
	testCases := []struct {
		url            string
		succeeds       bool
		expectedResult string
	}{
		{
			"http://golang.org",
			true,
			"http://golang.org a4dbbd512dcc9d4a9d0be4f36d78d216",
		}, {
			"http://duckduckgo.com",
			true,
			"http://duckduckgo.com 05d0fe2342405ea818f6a44511ac53f5",
		}, {
			"http://www.globo.com",
			false,
			`http://www.globo.com (Error: Get "http://www.globo.com": some error)`,
		}, {
			"https://g1.globo.com",
			false,
			`https://g1.globo.com (Error: Get "https://g1.globo.com": some error)`,
		},
	}
	for _, testCase := range testCases {
		client := newClientOfMineWithHTTPClient(&http.Client{
			Transport: roundTripper(roundTripperOfSuccess),
		})
		if !testCase.succeeds {
			client.Transport = roundTripper(roundTripperOfFailure)
		}
		pageResult := client.doGetPage(testCase.url)
		if pageResult.String() != testCase.expectedResult {
			t.Errorf("Got %s but expected %s", pageResult, testCase.expectedResult)
		}
	}
}

func ExampleRequest() {
	testCases := []struct {
		urls []string
	}{
		{
			[]string{
				"https://www.formula1.com",
				"https://www.mlb.com",
				"http://coffeegeek.com",
				"https://pabloaguiar.me",
			},
		}, {
			[]string{
				"https://github.com",
				"https://gitlab.com",
				"https://hub.docker.com",
				"https://www.manning.com",
				"https://www.mlb.com",
				"https://gitlab.com",
				"https://hub.docker.com",
				"http://coffeegeek.com",
			},
		}, {
			[]string{},
		},
	}
	client := newClientOfMineWithHTTPClient(&http.Client{
		Transport: roundTripper(roundTripperOfSuccess),
	})
	for _, testCase := range testCases {
		client.request(testCase.urls, 1)
	}
	// Output:
	// https://www.formula1.com a6bbc57376114596ddf51c4f9f7fdcb0
	// https://www.mlb.com b6e7c2d493778a4c7662273355008762
	// http://coffeegeek.com d75c56b4b08d69152f7bdd1b84e94abc
	// https://pabloaguiar.me 2f4e55faf5ff25b6cd61d7d683375045
	// https://github.com 742b2de0628c6160ab72f79e0588d3c0
	// https://gitlab.com 64c178e48724f10785e70f47584d7fe4
	// https://hub.docker.com aac7dd2fdb127a6a30cdfd3d8ccab131
	// https://www.manning.com 499fd2ea7522318c1515c5a3c2572589
	// https://www.mlb.com b6e7c2d493778a4c7662273355008762
	// https://gitlab.com 64c178e48724f10785e70f47584d7fe4
	// https://hub.docker.com aac7dd2fdb127a6a30cdfd3d8ccab131
	// http://coffeegeek.com d75c56b4b08d69152f7bdd1b84e94abc
}
