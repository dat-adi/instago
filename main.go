package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
    "github.com/dat-adi/instago/model/user"
    "github.com/dat-adi/instago/model/post"
)

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage Endpoint Hit...\n")
	fmt.Fprintf(response, "Please check the other endpoints now.")
}

/*
   Routing for the requests is done here
*/
func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users/{id}", getUser)
	http.HandleFunc("/users", postUser)
	http.HandleFunc("/posts/{id}", getPost)
	http.HandleFunc("/posts", postPost)
	http.HandleFunc("/posts/users/{id}", getAllPosts)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

var client *mongo.Client

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/insta"))
	handleRequests()
}
