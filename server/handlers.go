package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/zlowram/godan/model"

	"github.com/jroimartin/monmq"
)

type MonmqCmd struct {
	Target  string
	Command string
}

func (s *server) tasksHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}]}\n")
		return
	}
	var task Task
	err = json.Unmarshal(data, &task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"code\":\"400\",\"title\":\"Bad Request\",\"detail\":\"Invalid json format.\"}]}\n")
		return
	}
	tm := newTaskManager(s.client, s.pm)
	go tm.runTasks(task)
	fmt.Fprintln(w, "{\"status\": \"success\"}")
}

func (s *server) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(s.supervisor.Status())
	if err != nil {
		log.Println("Error marshaling:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}]}\n")
		return
	}
	fmt.Fprintf(w, "%s", b)
}

func (s *server) setStatusHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}]}\n")
		return
	}
	var cmd MonmqCmd
	if err = json.Unmarshal(data, &cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"code\":\"400\",\"title\":\"Bad Request\",\"detail\":\"Invalid json format.\"}]}\n")
		return
	}
	var command monmq.Command
	switch {
	case cmd.Command == "softshutdown":
		command = monmq.SoftShutdown
	case cmd.Command == "hardshutdown":
		command = monmq.HardShutdown
	case cmd.Command == "pause":
		command = monmq.Pause
	case cmd.Command == "resume":
		command = monmq.Resume
	}
	if err = s.supervisor.Invoke(command, cmd.Target); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"code\":\"400\",\"title\":\"Bad Request\",\"detail\":\"Invalid command for target.\"}]}\n")
		return
	}
}

func (s *server) queryHandler(w http.ResponseWriter, r *http.Request) {
	f := extractFilters(r)

	result, err := s.pm.QueryBanners(f)
	if err != nil {
		log.Println("Error querying banners:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}]}\n")
		return
	}
	jsoned, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshaling:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}]}\n")
		return
	}
	fmt.Fprintf(w, string(jsoned))
}

func extractFilters(request *http.Request) model.Filters {
	var p, s []string
	var r string

	values := request.URL.Query()
	ip := values.Get(":ip")
	if ip != "" {
		ip = strings.TrimPrefix(ip, "/")
	}
	if values["port"] != nil {
		p = strings.Split(values["port"][0], ",")
	}
	if values["service"] != nil {
		s = strings.Split(values["service"][0], ",")
	}
	if values["regexp"] != nil {
		r = values["regexp"][0]
	}
	filters := model.Filters{
		Ip:       ip,
		Ports:    p,
		Services: s,
		Regexp:   r,
	}
	return filters
}
