# Warning! This repo was no longer being developed!
I recommand using [goky](https://github.com/philippgille/gokv) instead

# keyv
Keyv provides a key-value interface to access database. By changing the adapter, you can change the database you use.  
Inspired by [lukechilds/keyv](https://github.com/lukechilds/keyv), a similar nodejs module.  

# Adapters
[sqlite3](https://github.com/simba-fs/keyvSqlite3)

# Usage
```go
package main

import (
	"fmt"
	"os"

	"github.com/simba-fs/keyv"
	_ "github.com/simba-fs/keyvSqlite3"
)

// checkErr return true if err != nil, and print err to stderr
func checkErr(err error) bool {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return true
	}
	return false
}

type user struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	// create a new connect with namesapce
	pw, err := keyv.New("sqlite3://database.sqlite", "password")
	if err != nil {
		fmt.Println(err)
		return
	}

	// create another connect with a different namesapce
	usr, err := keyv.New("sqlite3://database.sqlite", "user")
	if err != nil {
		fmt.Println(err)
		return
	}

	// write and read a string
	checkErr(pw.Set("peter", "P@ssw0Rd"))
	fmt.Println(pw.GetString("peter")) // P@ssw0Rd <nil>

	// write a struct
	sean := user{
		Username: "Sean",
		Email:    "seam@example.com",
	}
	checkErr(usr.Set("sean", sean))

	// read a struct
	u := user{}
	checkErr(usr.Get("sean", &u))
	fmt.Printf("%#v\n", u) // main.user{Username:"Sean", Email:"seam@example.com"}

	// list keys in namesapce
	fmt.Println(usr.Keys()) // [sean] <nil>
	fmt.Println(pw.Keys()) // [peter] <nil>
}
```
