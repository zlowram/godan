package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/jroimartin/monmq"
	"github.com/jroimartin/orujo"
)

type MonmqCmd struct {
	Target  string
	Command string
}

type Filters struct {
	Ip       string
	Ports    []string
	Services []string
	Regexp   string
}

type Banner struct {
	Ip      string
	Port    string
	Service string
	Content string
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
	switch {
	case r.Method == "POST":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			orujo.RegisterError(w, fmt.Errorf("Reading POST:", err))
			return
		}
		var cmd MonmqCmd
		if err = json.Unmarshal(data, &cmd); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			orujo.RegisterError(w, fmt.Errorf("Parsing json:", err))
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
			orujo.RegisterError(w, fmt.Errorf("Invoking monmq command:", err))
			return
		}
	case r.Method == "GET":
		b, err := json.Marshal(s.supervisor.Status())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			orujo.RegisterError(w, fmt.Errorf("error marshaling status:", err))
			return
		}
		fmt.Fprintf(w, "%s", b)
	}
}

func (s *server) queryHandler(w http.ResponseWriter, r *http.Request) {
	f := extractFilters(r)

	t := template.Must(template.New("query").Parse(queryTemplate))
	query := &bytes.Buffer{}
	err := t.Execute(query, f)
	if err != nil {
		log.Println("executing template:", err)
	}

	stmt, err := s.database.Prepare(query.String())
	if err != nil {
		log.Println("preparing statement:", err)
	}

	data := make([]interface{}, 0, len(f.Ports))
	if f.Ip != "" {
		data = append(data, f.Ip)
	}
	for _, v := range f.Ports {
		data = append(data, v)
	}
	for _, v := range f.Services {
		data = append(data, v)
	}
	if f.Regexp != "" {
		data = append(data, f.Regexp)
	}

	rows, err := stmt.Query(data...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		orujo.RegisterError(w, fmt.Errorf("error querying:", err))
		return
	}
	defer rows.Close()

	var result []Banner
	for rows.Next() {
		var curr Banner
		if err := rows.Scan(&curr.Ip, &curr.Port, &curr.Service, &curr.Content); err != nil {
			log.Fatal(err)
		}
		result = append(result, curr)
	}
	jsoned, err := json.Marshal(result)
	if err != nil {
		log.Println("marshaling:", err)
	}
	fmt.Fprintf(w, string(jsoned))
}

const (
	queryTemplate = `SELECT DISTINCT INET_NTOA(ip), port, service, content FROM banners
{{if or .Ip .Ports .Services .Regexp}} WHERE{{end}}
{{if .Ip}} ip = INET_ATON((?)){{end}}
{{if and .Ip .Ports}} AND{{end}}{{if .Ports}} port IN ({{range $i, $v := .Ports}}{{if $i}},{{end}}?{{end}}){{end}}
{{if and .Services (or .Ports .Services)}} AND{{end}}{{if .Services}} service IN ({{range $i, $v := .Services}}{{if $i}},{{end}}?{{end}}){{end}}
{{if and .Regexp (or .Ip .Ports .Services)}} AND{{end}}{{if .Regexp}} content regexp (?){{end}};`
)

func extractFilters(request *http.Request) Filters {
	var p, s []string
	var r string

	values := request.URL.Query()
	ip := strings.TrimPrefix(request.URL.Path, "/ips")
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
	filters := Filters{
		Ip:       ip,
		Ports:    p,
		Services: s,
		Regexp:   r,
	}
	return filters
}
