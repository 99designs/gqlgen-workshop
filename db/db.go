package db

import "strconv"

type User struct {
	ID    string
	Name  string
	Likes []string
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
		Name: "Billy",
	},
	"6": &User{
		ID:   "6",
		Name: "Cameron",
	},
}

var nextID = 7

func GetUser(id string) *User {
	return users[id]
}

func AddUser(name string) {
	id := strconv.Itoa(nextID)
	users[id] = &User{
		ID:   id,
		Name: name,
	}
	nextID++
}

func (u *User) Like(movieID string) {
	for _, id := range u.Likes {
		if movieID == id {
			return
		}
	}
	u.Likes = append(u.Likes, movieID)
}
