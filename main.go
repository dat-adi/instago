package main

import (
    "context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage Endpoint Hit.\n")
    fmt.Fprintf(response, "Please check the other endpoints now.")
}

/*
   Routing for the requests is done here
*/
func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":9000", nil))
}


/*
   From here, consider splitting it up into separate files.
   The one below this comment is model/user.go
*/


func main() {
    fmt.Println("Starting the application...")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/insta"))
	handleRequests()
}

/*
   From here, below this is the model/post.go
*/

