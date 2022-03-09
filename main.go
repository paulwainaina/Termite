package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter ,r* http.Request ){
	io.WriteString(w,"Welcome To termitarium")
}

func main(){
	fmt.Println("Welcome")

	r:=mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	r.HandleFunc("/",indexHandler).Methods("GET").Schemes("https")
	err:=http.ListenAndServeTLS(":8443","server.crt","server.key",r)
	if err!=nil{
		fmt.Println(err)
	}
}