package db

import (
	"strconv"
	"sync"
)

type User struct {
	ID    string
	Name  string
	Likes []string
	m     sync.Mutex
}

var users = map[string]*User{
	"1": &User{
		ID:   "1",
		Name: "Chris",
	},
	"2": &User{
		ID:   "2",
		Name: "Steph",
	},
	"3": &User{
		ID:   "3",
		Name: "Peter",
	},
	"4": &User{
		ID:   "4",
		Name: "Tas",
	},
	"5": &User{
		ID:   "5",
		Name: "Cameron",
	},
	"6": &User{
		ID:   "6",
		Name: "Tim",
	},
}

var nextID = struct {
	id int
	m  sync.Mutex
}{7, sync.Mutex{}}

func GetUser(id string) *User {
	u, ok := users[id]
	if !ok {
		return nil
	}
	return u
}

func AddUser(name string) {
	nextID.m.Lock()
	defer nextID.m.Unlock()
	id := strconv.Itoa(nextID.id)
	users[id] = &User{
		ID:   id,
		Name: name,
	}
	nextID.id++
}

func (u *User) Like(movieID string) {
	u.m.Lock()
	defer u.m.Unlock()
	for _, id := range u.Likes {
		if movieID == id {
			return
		}
	}
	u.Likes = append(u.Likes, movieID)
}
