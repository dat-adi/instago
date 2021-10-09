package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
    user "github.com/dat-adi/instago/model/user"
    post "github.com/dat-adi/instago/model/post"
)

// Defining a simple homepage to land on
func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage Endpoint Hit...\n")
	fmt.Fprintf(response, "Please check the other endpoints now.")
    fmt.Println("From the homepage => ", request.URL)
}

/*
   Routing for the requests is done here
*/
func handleRequests() {
    // Routes to the various pages
	http.HandleFunc("/", homePage)
    http.HandleFunc("/users/:id", user.getUser)
	http.HandleFunc("/users", user.postUser)
    http.HandleFunc("/posts/:id", post.getPost)
	http.HandleFunc("/posts", post.postPost)
    http.HandleFunc("/posts/users/:id", post.getPostsByUserId)

//	http.HandleFunc("/users/all", user.getAllUsers)

    // Starting up the server
	log.Fatal(http.ListenAndServe(":9000", nil))
}

var client *mongo.Client

// The main function, where everything starts to take place
func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/insta"))

    // Calling a method for routing
	handleRequests()
}
