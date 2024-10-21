package authorization

import (
	"WebServer/internal/config/app"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"
)

var buffer = []string{}
var clearingStarted = false

func startClearing(m *sync.Mutex) {
	if clearingStarted {
		return
	}
	clearingStarted = true
	for {
		<-time.After(time.Hour * 24)
		m.Lock()
		buffer = nil
		m.Unlock()
		log.Println("Buffer cleared")
	}
}

// [ ] Auth Service
func Authorize(w http.ResponseWriter, r *http.Request) bool {
	log.Println("Check authorization for", r.RemoteAddr)
	m := sync.Mutex{}
	if slices.Contains(buffer, r.RemoteAddr) {
		return true
	}
	if len(buffer) > 1000 {
		log.Println("Clearing buffer")
		buffer = buffer[900:]
	}
	go startClearing(&m)

	login, err := r.Cookie("login")
	if err != nil || login.Value != app.LOGIN {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return false
	}

	password, err := r.Cookie("password")
	if err != nil || password.Value != app.PASSWORD {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return false
	}
	buffer = append(buffer, r.RemoteAddr)
	return true
}
