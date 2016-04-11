package main

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	UserId   string
	Username string
	Email    string
	Hash     string
	Role     string
}
