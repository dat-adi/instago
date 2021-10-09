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
    "github.com/dat-adi/instago/routes/route"
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
    // Set up a router
    app := route.Router()

    // Routes to the various pages
	app.HandleFunc("^/$", homePage)
    app.HandleFunc("/users/([a-zA-Z0-9]+)$", user.getUser)
	app.HandleFunc("/users", user.postUser)
    app.HandleFunc("/posts/([a-zA-Z0-9]+)$", post.getPost)
	app.HandleFunc("/posts", post.postPost)
    app.HandleFunc("/posts/users/([a-zA-Z0-9]+)$", post.getPostsByUserId)

//	http.HandleFunc("/users/all", user.getAllUsers)

    // Starting up the server
    err := http.ListenAndServe(":9000", app)
    if err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
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
