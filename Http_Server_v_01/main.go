package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	TraceLog *log.Logger
	LogFile  os.File
	// errStr   string

	LOG_INFO_NAME string = "log.info.txt"
	LOG_INFO_ERR  string = "log.err.txt"
	LOG_FILENAME  string = "app.log"
)

const (
	TraceApp = true
)

type ProtocolType struct {
	ProtocolID string `json:"ProtocolID"`
	Comment    string `json:"Comment"`
}

var ProtocolList []ProtocolType

func handleRequests() {

	// new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// main
	myRouter.HandleFunc("/", homePage)
	// myRouter.HandleFunc("/s", StopServer)
	// Protocol
	myRouter.HandleFunc("/Protocols", returnAllProtocols).Methods("GET")
	myRouter.HandleFunc("/Protocols", PostProtocols).Methods("POST")

	myRouter.HandleFunc("/Protocol/{id}", returnOneProtocol)
	myRouter.HandleFunc("/Protocol", createNewProtocol).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", myRouter))

	// Стандартный (http)
	// http.HandleFunc("/", homePage)
	// http.HandleFunc("/Protocols", returnAllProtocols)
	// log.Fatal(http.ListenAndServe(":8081", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")

	fmt.Println("Endpoint Hit: homePage")
}

// func StopServer(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Server ShutDown now !")

// 	// http.Server.Handler.ShutDown()
// 	var srv http.Server

// 	idleConnsClosed := make(chan struct{})

// 	go func() {
// 		sigint := make(chan os.Signal, 1)
// 		signal.Notify(sigint, os.Interrupt)
// 		<-sigint

// 		// We received an interrupt signal, shut down.
// 		if err := srv.Shutdown(context.Background()); err != nil {
// 			// Error from closing listeners, or context timeout:
// 			log.Printf("HTTP server Shutdown: %v", err)
// 		}
// 		close(idleConnsClosed)
// 	}()

// 	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
// 		// Error starting or closing listener:
// 		log.Fatalf("HTTP server ListenAndServe: %v", err)
// 	}

// 	<-idleConnsClosed
// }

func returnAllProtocols(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllProtocols")

	json.NewEncoder(w).Encode(ProtocolList)
}

// GET http://localhost:8081/Protocols

func PostProtocols(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TEST POST!")
}

func returnOneProtocol(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, Pr := range ProtocolList {
		if Pr.ProtocolID == key {
			json.NewEncoder(w).Encode(Pr)
		}
	}

	// Просто вывод в браузер
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Key: "+key)
}

func createNewProtocol(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))
} // createNewProtocol

func main() {

	LogFile, err := os.OpenFile(LOG_FILENAME, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}
	defer LogFile.Close()

	// Трассировка
	if TraceApp {
		TraceLog := log.New(LogFile, "Trace\t", log.LstdFlags)
		TraceLog.SetFlags(log.LstdFlags)
		TraceLog.Println("Start programm...")
	}

	ProtocolList = []ProtocolType{
		{ProtocolID: "1", Comment: "Comm1"},
		{ProtocolID: "2", Comment: "Comm2"},
	}

	handleRequests()
}
