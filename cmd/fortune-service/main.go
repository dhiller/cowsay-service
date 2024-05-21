/*
 * MIT License
 *
 * Copyright (c) 2024 Daniel Hiller
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync/atomic"
)

var (
	addr              string
	cowsayServiceAddr string
	cowsays           []string
	cowsayCounter     atomic.Int32
)

func main() {
	flag.StringVar(&addr, "address", ":9090", "the address of the server (host and port)")
	flag.StringVar(&cowsayServiceAddr, "cowsay-service-address", "http://localhost:8080", "the address of the cowsay server (host and port)")
	flag.Parse()
	http.Handle("/", http.HandlerFunc(getFortune))
	log.Printf("fortune-service starting")
	var err error
	cowsays, err = fetchCowsays()
	if err != nil {
		log.Fatalf("fortune-service failure: %v", err)
	}
	log.Printf("cowsays: %v", cowsays)
	if err = http.ListenAndServe(addr, nil); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("fortune-service failure: %v", err)
	}
}

func getFortune(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling request %s from %s", req.URL, req.RemoteAddr)
	fortuneOutput, err := exec.Command("fortune").Output()
	if err != nil {
		log.Printf("ERROR: fortune-service - fortune output retrieval: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var currentCounterValue int32
	var swapped bool
	for swapped != true {
		currentCounterValue = cowsayCounter.Load()
		nextCounterValue := currentCounterValue + 1
		swapped = cowsayCounter.CompareAndSwap(currentCounterValue, nextCounterValue)
	}
	cowsayIndex := int(currentCounterValue) % len(cowsays)
	resp, err := http.Post(fmt.Sprintf("%s/cowsay/%s", cowsayServiceAddr, cowsays[cowsayIndex]), http.DetectContentType(fortuneOutput), bytes.NewReader(fortuneOutput))
	if err != nil {
		log.Printf("ERROR fortune-service - post cowsay service: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch resp.StatusCode {
	case http.StatusOK:
		break
	default:
		log.Printf("ERROR fortune-service - post response status %d", resp.StatusCode)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	_, err = w.Write(bodyBytes)
	if err != nil {
		log.Printf("ERROR fortune-service - write output: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func fetchCowsays() ([]string, error) {
	var result []string
	resp, err := http.Get(fmt.Sprintf("%s/cowsays", cowsayServiceAddr))
	if err != nil {
		return nil, fmt.Errorf("fortune-service: can't get cowsay variants: %v", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		break
	default:
		return nil, fmt.Errorf("fortune-service: failed to get cowsay variants: %v", err)
	}
	defer resp.Body.Close()
	cowsaysBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fortune-service: can't decode cowsay variants body: %v", err)
	}
	lines := strings.Split(string(cowsaysBody), "\n")
	for _, line := range lines[1:] {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		result = append(result, strings.Split(trimmedLine, " ")...)
	}
	return result, nil
}
