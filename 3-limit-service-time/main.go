//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mut sync.Mutex
}

const timeLimitSeconds = 10

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}
	u.mut.Lock()
	defer u.mut.Unlock()

	ch := make(chan bool)

	go func() {
		process()
		ch <- true
	}()

	select {
		case <-time.After(time.Second*timeLimitSeconds):
			return false
		case <-ch:
			return true
	}
}

func main() {
	RunMockServer()
}
