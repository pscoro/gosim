package main

import "log"

type user struct {
	id    uint32
	name  string
	email string
}

func newUser(name string, email string) (*user, error) {
	u := user{
		id:    0,
		name:  name,
		email: email,
	}
	db, err := newDatabaseConnection()
	defer db.closeDatabaseConnection()
	if err != nil {
		log.Panicln("Error connecting to database")
	}
	db.registerUser(&u)
	return &u, nil
}
