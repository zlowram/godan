package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"

	"gopkg.in/mgo.v2/bson"
)

func (s *server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}\n")
		return
	}
	c := s.db.DB("test").C("users")
	result := []User{}
	err := c.Find(bson.M{}).Select(bson.M{"_id": 0, "password": 0}).All(&result)
	if err != nil {
		fmt.Fprintln(w, "{}")
		return
	}
	ret, _ := json.Marshal(result)
	fmt.Fprintln(w, string(ret))
	return
}

func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	username := r.URL.Query().Get(":username")
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" && user["username"] != username {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}\n")
		return
	}

	c := s.db.DB("test").C("users")
	result := User{}
	err := c.Find(bson.M{"username": username}).Select(bson.M{"_id": 0, "password": 0, "role": 0}).One(&result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "{\"code\":\"404\",\"title\":\"Not Found\",\"detail\":\"User not found.\"}\n")
		return
	}
	ret, _ := json.Marshal(result)
	fmt.Fprintln(w, string(ret))
	return
}

func (s *server) newUserHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}\n")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newUser User
	err := decoder.Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"code\":\"400\",\"title\":\"Bad Request\",\"detail\":\"Invalid json format.\"}\n")
		return
	}

	h := sha256.New()
	h.Write([]byte(newUser.Password))
	newUser.Password = hex.EncodeToString(h.Sum(nil))

	c := s.db.DB("test").C("users")
	queryResult := User{}
	err = c.Find(bson.M{"username": newUser.Username}).One(&queryResult)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"code\":\"400\",\"title\":\"Bad Request\",\"detail\":\"User already exists.\"}\n")
		return
	}
	err = c.Insert(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}\n")
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/users/"+newUser.Username)
	return
}

func (s *server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	username := r.URL.Query().Get(":username")
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" && user["username"] != username {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}\n")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var updateUser User
	err := decoder.Decode(&updateUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"code\":\"400\",\"title\":\"Bad Request\",\"detail\":\"invalid json format.\"}\n")
		return
	}
	c := s.db.DB("test").C("users")
	current := User{}
	err = c.Find(bson.M{"username": username}).One(&current)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "{\"code\":\"404\",\"title\":\"Not Found\",\"detail\":\"User not found.\"}\n")
		return
	}

	if updateUser.Username != "" && user["role"] == "admin" {
		current.Username = updateUser.Username
	}
	if updateUser.Email != "" {
		current.Email = updateUser.Email
	}
	if updateUser.Password != "" {
		h := sha256.New()
		h.Write([]byte(updateUser.Password))
		current.Password = hex.EncodeToString(h.Sum(nil))
	}
	if updateUser.Role != "" && user["role"] == "admin" {
		current.Role = updateUser.Role
	}
	userQuery := bson.M{"username": username}
	err = c.Update(userQuery, current)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}\n")
		return
	}
	return
}

func (s *server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Auhtorization check
	user := context.Get(r, "user").(map[string]string)
	if user["role"] != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Access not authorized.\"}\n")
		return
	}

	username := r.URL.Query().Get(":username")
	c := s.db.DB("test").C("users")
	err := c.Remove(bson.M{"username": username})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}\n")
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}

func (s *server) loginHandler(w http.ResponseWriter, r *http.Request) {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var newLogin Login
	err := decoder.Decode(&newLogin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"code\":\"400\",\"title\":\"Bad Request\",\"detail\":\"Bad request.\"}\n")
		return
	}
	h := sha256.New()
	h.Write([]byte(newLogin.Password))
	hashedPasswd := hex.EncodeToString(h.Sum(nil))
	c := s.db.DB("test").C("users")
	queryResult := User{}
	err = c.Find(bson.M{"$and": []bson.M{bson.M{"username": newLogin.Username}, bson.M{"password": hashedPasswd}}}).One(&queryResult)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"code\":\"401\",\"title\":\"Unauthorized\",\"detail\":\"Login not valid.\"}\n")
		return
	}

	claims := make(map[string]string)
	claims["username"] = queryResult.Username
	claims["role"] = queryResult.Role
	token, err := s.auth.NewToken(claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"code\":\"500\",\"title\":\"Internal Server Error\",\"detail\":\"Something went wrong.\"}\n")
		return
	}
	fmt.Fprintf(w, "{\"username\": \"%s\", \"role\":\"%s\", \"accesToken\": \"%s\"}\n", queryResult.Username, queryResult.Role, token)
	return
}
