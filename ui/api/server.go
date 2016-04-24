package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/husobee/vestigo"
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
	m := vestigo.NewRouter()
	logger := log.New(os.Stdout, "[GODAN-UI] ", log.LstdFlags)
	logHandler := olog.NewLogHandler(logger, logLine)

	m.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:      []string{"*"},
		AllowCredentials: true,
		MaxAge:           3600 * time.Second,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	})

	m.Get("/users/:id", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.getUserHandler)).ServeHTTP,
	)
	m.Put("/users/:id", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.updateUserHandler)).ServeHTTP,
	)
	m.Delete("/users/:id", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.deleteUserHandler)).ServeHTTP,
	)
	m.Get("/users", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.getUsersHandler)).ServeHTTP,
	)
	m.Post("/users", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.newUserHandler)).ServeHTTP,
	)
	m.Post("/login", orujo.NewPipe(
		orujo.M(logHandler),
		http.HandlerFunc(s.loginHandler)).ServeHTTP,
	)

	m.Post("/tasks", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.newTaskHandler)).ServeHTTP,
	)
	m.Get("/ips/:ip", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.getIpHandler)).ServeHTTP,
	)
	m.Get("/ips", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.getIpHandler)).ServeHTTP,
	)
	m.Get("/status", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.getStatusHandler)).ServeHTTP,
	)
	m.Post("/status", orujo.NewPipe(
		orujo.M(logHandler),
		//s.auth,
		http.HandlerFunc(s.setStatusHandler)).ServeHTTP,
	)

	//http.Handle("/", m)
	fmt.Println("Listening on " + s.config.Local.Host + ":" + s.config.Local.Port + "...")
	log.Fatalln(http.ListenAndServe(s.config.Local.Host+":"+s.config.Local.Port, m))

	return nil
}

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}`
