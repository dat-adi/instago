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

/*
   From here, below this is the model/post.go
*/

type Post struct {
	ID       primitive.ObjectID `json:"id"`
	Caption  string             `json:"caption"`
	ImageURL string             `json:"image_url"`
	PostedAt *time.Time         `json:"posted_at"`
}

type Posts []Post

func getPost(response http.ResponseWriter, request *http.Request) {
	/*
	   return type => (post *Post)
	   parameter type => (id int64)
	   Get one post by ID
	*/
	if request.Method == "GET" {
		response.Header().Add("content-type", "application/json")
		params := request.URL.Query().Get("id")
		fmt.Fprintf(response, params)
		id, _ := primitive.ObjectIDFromHex(params)
		var post Post

		collection := client.Database("insta").Collection("post")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		// BROKEN for some reason
		err := collection.FindOne(ctx, User{ID: id}).Decode(&post)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
		json.NewEncoder(response).Encode(post)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
	fmt.Fprintf(response, "Endpoint Hit: Get Post by ID endpoint")
}

func postPost(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (post *Post)
	   Create a post
	*/
	if request.Method == "POST" {
		response.Header().Add("content-type", "application/json")
		var post Post
		json.NewDecoder(request.Body).Decode(&post)
		collection := client.Database("insta").Collection("post")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		result, _ := collection.InsertOne(ctx, post)
		json.NewEncoder(response).Encode(result)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}

	fmt.Fprintf(response, "Endpoint Hit: Create a Post endpoint")
}

func getAllPosts(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (user *User)
	   return type => array of (post *Post)?
	   Get all the posts made by a user
	*/
	if request.Method == "GET" {
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	} fmt.Fprintf(response, "Endpoint Hit: Get all user Posts endpoint")
}
