package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jroimartin/monmq"
	"github.com/jroimartin/orujo"
	olog "github.com/jroimartin/orujo-handlers/log"
	"github.com/jroimartin/rpcmq"
)

type server struct {
	config     Config
	logger     *log.Logger
	client     *rpcmq.Client
	supervisor *monmq.Supervisor
	database   *sql.DB
}

func newServer(cfg Config) *server {
	s := &server{
		config: cfg,
		logger: log.New(os.Stdout, "[godan] ", log.LstdFlags),
	}

	return s
}

func (s *server) start() error {
	s.client = rpcmq.NewClient("amqp://"+s.config.Rpcmq.Host+":"+s.config.Rpcmq.Port, s.config.Rpcmq.MsgQueue, s.config.Rpcmq.ReplyQueue, s.config.Rpcmq.Exchange, s.config.Rpcmq.ExchangeType)
	err := s.client.Init()
	if err != nil {
		log.Fatalf("Init rpcmq: %v", err)
	}
	defer s.client.Shutdown()

	s.supervisor = monmq.NewSupervisor("amqp://10.0.1.3:5672",
		"mon-replies", "mon-exchange")
	if err := s.supervisor.Init(); err != nil {
		log.Fatalf("Init monmq: %v", err)
	}
	defer s.supervisor.Shutdown()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", s.config.DB.Username, s.config.DB.Password, s.config.DB.Host, s.config.DB.Port, s.config.DB.Name)
	s.database, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer s.database.Close()

	_, err = s.database.Exec("CREATE TABLE IF NOT EXISTS banners (ip INT UNSIGNED, port INT UNSIGNED, service VARCHAR(50), content TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	websrv := orujo.NewServer("localhost:8080")
	fmt.Println("Listening on localhost:8080...")

	logHandler := olog.NewLogHandler(s.logger, logLine)

	websrv.Route(`^/tasks$`,
		http.HandlerFunc(s.tasksHandler),
		orujo.M(logHandler)).Methods("POST")

	websrv.Route(`^/status$`,
		http.HandlerFunc(s.statusHandler),
		orujo.M(logHandler)).Methods("GET", "POST")

	websrv.Route(`^/ips\??$`,
		http.HandlerFunc(s.queryHandler),
		orujo.M(logHandler)).Methods("GET")

	websrv.Route(`^/ips/(?:\d{1,3}\.){3}\d{1,3}\??$`,
		http.HandlerFunc(s.queryHandler),
		orujo.M(logHandler)).Methods("GET")

	log.Fatalln(websrv.ListenAndServe())

	return nil
}

const (
	logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
		{{range  $err := .Errors}}  Err: {{$err}}
		{{end}}`
	errLine = `{"error":"{{.}}"}`
)
