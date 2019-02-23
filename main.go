package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func main() {
	// read configuration
	config, err := ParseConfig("config.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("setting root path to %s\n", config.RootPath)

	// initialize in-memory session store
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	store = sessions.NewCookieStore(randBytes)

	// router
	r := mux.NewRouter()

	// some routes appear here...

	// receive signal
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go http.ListenAndServe(":10101", r)
	<-sig

	fmt.Printf("Signal received. Exit...\n")
	close(sig)
}
