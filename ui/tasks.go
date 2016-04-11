package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

func (s *server) newTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}]}\n")
		return
	}

	forwardPost(w, r, "http://"+s.config.Godan.Host+":"+s.config.Godan.Port+r.URL.RequestURI())
	return
}

func (s *server) getIpHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" && user["role"] != "user" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}]}\n")
		return
	}

	forwardGet(w, r, "http://"+s.config.Godan.Host+":"+s.config.Godan.Port+r.URL.RequestURI())
	return
}

func (s *server) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}]}\n")
		return
	}

	forwardGet(w, r, "http://"+s.config.Godan.Host+":"+s.config.Godan.Port+r.URL.RequestURI())
	return
}

func (s *server) setStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}]}\n")
		return
	}

	forwardPost(w, r, "http://"+s.config.Godan.Host+":"+s.config.Godan.Port+r.URL.RequestURI())
	return
}

func forwardGet(w http.ResponseWriter, r *http.Request, dest string) {
	req, err := http.NewRequest("GET", dest, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(resp.StatusCode)
	fmt.Fprintf(w, string(body))
	return
}

func forwardPost(w http.ResponseWriter, r *http.Request, dest string) {
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", dest, bytes.NewBuffer(rbody))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(resp.StatusCode)
	fmt.Fprintf(w, string(body))
	return
}
