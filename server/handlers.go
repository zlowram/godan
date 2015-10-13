package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jroimartin/orujo"
)

type Filters struct {
	ports    []string
	services []string
	regexp   string
}

func (s *server) tasksHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		orujo.RegisterError(w, fmt.Errorf("Reading POST:", err))
		return
	}
	tasks := []string{string(data)}
	tm := newTaskManager(s.client)
	go tm.runTasks(tasks)
	fmt.Fprintln(w, "{\"status\": \"success\"}")
}

func (s *server) statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Current status requested.")
}

func (s *server) allIpsHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	f := extractFilters(values)
	fmt.Fprintln(w, "Ports:", f.ports)
	fmt.Fprintln(w, "Services:", f.services)
	fmt.Fprintln(w, "Regexp:", f.regexp)
}

func (s *server) ipsHandler(w http.ResponseWriter, r *http.Request) {
	ip := strings.TrimPrefix(r.URL.Path, "/ips/")
	values := r.URL.Query()
	f := extractFilters(values)
	fmt.Fprintln(w, "IP:", ip)
	fmt.Fprintln(w, "Ports:", f.ports)
	fmt.Fprintln(w, "Services:", f.services)
	fmt.Fprintln(w, "Regexp:", f.regexp)
}

func extractFilters(values url.Values) Filters {
	var p, s []string
	var r string
	if values["port"] != nil {
		p = strings.Split(values["port"][0], ",")
	}
	if values["service"] != nil {
		s = strings.Split(values["service"][0], ",")
	}
	if values["regexp"] != nil {
		r = values["regexp"][0]
	}
	filters := Filters{
		ports:    p,
		services: s,
		regexp:   r,
	}
	return filters
}
