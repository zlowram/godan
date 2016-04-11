package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bmizerany/pat"
	"github.com/jroimartin/orujo"
	olog "github.com/jroimartin/orujo/log"
	"github.com/zlowram/zmiddlewares/jwtauth"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type server struct {
	config Config
	db     *mgo.Session
	auth   *jwtauth.AuthHandler
}

func newServer(config Config) *server {
	var err error
	db, err := mgo.Dial(config.DB.Host + ":" + config.DB.Port)
	if err != nil {
		log.Fatal("ERROR connecting to DB:", err)
	}

	queryResult := User{}
	c := db.DB("test").C("users")
	err = c.Find(bson.M{}).One(&queryResult)
	if err != nil {
		defaultUser := User{
			UserId:   "0",
			Username: "admin",
			Email:    "admin@localhost.com",
			Role:     "admin",
		}
		h := sha256.New()
		h.Write([]byte("admin"))
		defaultUser.Hash = hex.EncodeToString(h.Sum(nil))
		err = c.Insert(defaultUser)
		if err != nil {
			log.Fatal("Error adding default user:", err)
		}
	}

	return &server{
		config: config,
		db:     db,
		auth:   jwtauth.NewAuthHandler(config.Local.PrivateKey, config.Local.PublicKey),
	}
}

func (s *server) start() error {
	m := pat.New()
	logger := log.New(os.Stdout, "[GODAN-UI] ", log.LstdFlags)
	logHandler := olog.NewLogHandler(logger, logLine)

	m.Get("/users/:id/", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.getUserHandler)),
	)
	m.Put("/users/:id/", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.updateUserHandler)),
	)
	m.Del("/users/:id/", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.deleteUserHandler)),
	)
	m.Get("/users/", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.getUsersHandler)),
	)
	m.Post("/users/", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.newUserHandler)),
	)
	m.Post("/login/", orujo.NewPipe(
		orujo.M(logHandler),
		http.HandlerFunc(s.loginHandler)),
	)

	m.Post("/tasks", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.newTaskHandler)),
	)
	m.Get("/ips/", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.getIpHandler)),
	)
	m.Get("/status", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.getStatusHandler)),
	)
	m.Post("/status", orujo.NewPipe(
		orujo.M(logHandler),
		s.auth,
		http.HandlerFunc(s.setStatusHandler)),
	)

	http.Handle("/", m)
	fmt.Println("Lisening on " + s.config.Local.Host + ":" + s.config.Local.Port + "...")
	log.Fatalln(http.ListenAndServe(s.config.Local.Host+":"+s.config.Local.Port, nil))

	return nil
}

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}`
