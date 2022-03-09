package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"encoding/json"

	"example.com/Db"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

var database Db.Database

type User struct{ 
	Name string;
	Age int;
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	io.WriteString(w, "Welcome To termitarium")
}
func getDataHandler(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{}
	result, err := database.QueryMany(os.Getenv("DB"), os.Getenv("COL"), filter)
	var users []User
	if err == nil {
		for _, x := range result {
			doc, _ := bson.Marshal(x)
			var user User
			bson.Unmarshal(doc, &user)
			users = append(users, user)
		}
	}
	json.NewEncoder(w).Encode(users)
}
func insertDataHandler(w http.ResponseWriter, r *http.Request) {
	pname := r.FormValue("pname")
	page,err :=strconv.Atoi(r.FormValue("page"))
	if err!=nil||len(pname)<=0|| page==0 {
		http.Error(w, "Missing fields", http.StatusForbidden)
		return
	}
	document := bson.M{"Age": page, "Name": pname}
	var result interface{}
	result, err1 := database.InsertOne(os.Getenv("DB"), os.Getenv("COL"), document)
	if err1 != nil {
		http.Error(w, "Insertion failed", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func main() {
	fmt.Println("Welcome")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading environment")
	}

	database.Connect()
	defer database.DisConnect()
	
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(os.Getenv("ASSETS")))))
	r.HandleFunc("/", indexHandler).Methods("GET").Schemes("https")
	r.HandleFunc("/insertData", insertDataHandler).Methods("POST").Schemes("https")
	r.HandleFunc("/getData", getDataHandler).Methods("GET").Schemes("https")

	err = http.ListenAndServeTLS(os.Getenv("PORT"), os.Getenv("CERT"), os.Getenv("PRIVKEY"), r)
	if err != nil {
		fmt.Println(err)
	}
}
