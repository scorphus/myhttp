// This file is part of myhttp

// Licensed under the BSD-3-Clause license:
// https://opensource.org/licenses/BSD-3-Clause
// Copyright (c) 2020, Pablo S. Blum de Aguiar

package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// defaultClient is the default http.Client used to issue requests, with
// DefaultTransport. It is safe for concurrent use by multiple goroutines
var defaultClient = &http.Client{Timeout: timeout()}

// myUserAgent holds the most common User Agent on the date of this writing
const myUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"

type httpClient struct{ *http.Client }

type pageResult struct {
	url  string
	data [16]byte
	err  error
}

func newClientOfMine() *httpClient {
	return newClientOfMineWithHTTPClient(defaultClient)
}

func newClientOfMineWithHTTPClient(c *http.Client) *httpClient {
	return &httpClient{c}
}

func timeout() time.Duration {
	if seconds, err := strconv.Atoi(os.Getenv("MYHTTP_TIMEOUT")); err == nil {
		return time.Second * time.Duration(seconds)
	}
	return time.Second * 30
}

func newPageResult(url string, data []byte, err error) *pageResult {
	return &pageResult{url, md5.Sum(data), err} // keep in memory only the necessary bits
}

// request prepares the requests to run at a pace of maxConcurrent
// at a time and pipes each pageResult into the pageFeed channel
func (client *httpClient) request(urls []string, maxConcurrent uint64) {
	wg := sync.WaitGroup{}
	wg.Add(len(urls))
	pacer := make(chan struct{}, maxConcurrent)
	for _, url := range urls {
		pacer <- struct{}{}
		go func(url string) {
			defer wg.Done()
			defer func() { <-pacer }()
			fmt.Println(client.doGetPage(url)) // TODO: print elsewhere (e.g. pipe it into a channel)
		}(url)
	}
	wg.Wait()
}

// doGetPage dispatches the request and converts the response into a pageResult
func (client *httpClient) doGetPage(url string) *pageResult {
	data, err := client.doGet(url)
	if err != nil {
		return newPageResult(url, nil, err)
	}
	return newPageResult(url, data, nil)
}

// doGet makes a `GET url` request
func (client *httpClient) doGet(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", myUserAgent) // request url as if a regular user
	req.Header.Set("Connection", "close")     // no need to keep connections open
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (page pageResult) String() string {
	if page.err != nil {
		return fmt.Sprintf("%s (Error: %s)", page.url, page.err.Error())
	}
	return fmt.Sprintf("%s %x", page.url, page.data)
}
