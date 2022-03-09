package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func indexHandler(w http.ResponseWriter ,r* http.Request ){
	fmt.Println(r.Header)
	io.WriteString(w,"Welcome To termitarium")
}

func main(){
	fmt.Println("Welcome")
	err:=godotenv.Load()
	if err!=nil{
		fmt.Println("Error loading environment")
	}

	r:=mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(os.Getenv("ASSETS")))))
	r.HandleFunc("/",indexHandler).Methods("GET").Schemes("https")
	err=http.ListenAndServeTLS(os.Getenv("PORT"),os.Getenv("CERT"),os.Getenv("PRIVKEY"),r)
	if err!=nil{
		fmt.Println(err)
	}
}