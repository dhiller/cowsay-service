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
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func getListOfCowsays(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling request %s from %s", req.URL, req.RemoteAddr)
	output, err := exec.Command("cowsay", "-l").Output()
	if err != nil {
		log.Printf("ERROR: cowsay: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(output)
	if err != nil {
		log.Printf("ERROR: cowsay: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getCowsay(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling request %s from %s", req.URL, req.RemoteAddr)
	//req.URL contains the cowsay we want
	pathElements := strings.Split(req.URL.Path, "/")
	args := []string{}
	if len(pathElements) > 2 {
		args = append(args, "-f", pathElements[2])
	}
	//req.Body contains what it should say
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("ERROR: cowsay: %v", err)
		w.WriteHeader(500)
		return
	}
	toSay := string(bodyBytes)
	args = append(args, toSay)
	output, err := exec.Command("cowsay", args...).Output()
	if err != nil {
		log.Printf("ERROR: cowsay: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(output)
	if err != nil {
		log.Printf("ERROR: cowsay: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var addr string

func main() {
	flag.StringVar(&addr, "address", ":8080", "the address of the server (host and port)")
	flag.Parse()
	http.Handle("/cowsay/", http.HandlerFunc(getCowsay))
	http.Handle("/cowsay", http.HandlerFunc(getCowsay))
	http.Handle("/cowsays", http.HandlerFunc(getListOfCowsays))
	log.Printf("cowsay-service starting")
	if err := http.ListenAndServe(addr, nil); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("cowsay-service failure: %v", err)
	}
}
