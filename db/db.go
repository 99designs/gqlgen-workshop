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

var users sync.Map

func init() {
	users.Store("1", &User{
		ID:   "1",
		Name: "Chris",
	})
	users.Store("2", &User{
		ID:   "2",
		Name: "Steph",
	})
	users.Store("3", &User{
		ID:   "3",
		Name: "Peter",
	})
	users.Store("4", &User{
		ID:   "4",
		Name: "Tas",
	})
	users.Store("5", &User{
		ID:   "5",
		Name: "Cameron",
	})
	users.Store("6", &User{
		ID:   "6",
		Name: "Tim",
	})

	nextID.id = 7
}

var nextID struct {
	id int
	m  sync.Mutex
}

func GetUser(id string) *User {
	u, ok := users.Load(id)
	if !ok {
		return nil
	}
	return u.(*User)
}

func AddUser(name string) {
	nextID.m.Lock()
	id := strconv.Itoa(nextID.id)
	users.Store(id, &User{
		ID:   id,
		Name: name,
	})
	nextID.id++
	nextID.m.Unlock()
}

func (u *User) Like(movieID string) {
	u.m.Lock()
	for _, id := range u.Likes {
		if movieID == id {
			return
		}
	}
	u.Likes = append(u.Likes, movieID)
	u.m.Unlock()
}
