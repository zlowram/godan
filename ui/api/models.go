package main

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Username string
	Email    string
	Password string
	Role     string
}
