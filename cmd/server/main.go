package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/NjiruClinton/reproxy-mobile-bff/initializations"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func initEnvFile(wg *sync.WaitGroup) {
	defer wg.Done()
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Couldn't load env file %s\n", err)
	}
}

func initEnvs(wg *sync.WaitGroup) {
	defer wg.Done()
	err := godotenv.Load(os.Getenv(initializations.APP_PROFILE) + ".env")
	if err != nil {
		fmt.Printf("could not load envsss %s\n", err)
	}
}

func initServer() {
	r := mux.NewRouter()

	r.HandleFunc("/login", handleLogin).Methods("POST")

	log.Fatal(http.ListenAndServe(":4100", r))
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go initEnvFile(wg)
	go initEnvs(wg)
	wg.Wait()
}
