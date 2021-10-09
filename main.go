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

/*
   From here, consider splitting it up into separate files.
   The one below this comment is model/user.go
*/

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"-,omitempty" bson:"-,omitempty"`
}

type Users []User

var client *mongo.Client

func getUser(response http.ResponseWriter, request *http.Request) {
	/*
	   return type => (user *User)
	   parameter type => (id int64)
	   Get one user by ID
	*/
	// Checking the HTTP Request Method
	if request.Method == "GET" {
		// Addition of a Header to the response
		response.Header().Add("content-type", "application/json")

		// Parsing the URL for parameters
		params := request.URL.Query().Get("id")
		fmt.Fprintf(response, params)

		// Converting the params into an ObjectID
		id, _ := primitive.ObjectIDFromHex(params)

		// Creating a variable user
		var user User

		// Connecting to MongoDB's user collection
		collection := client.Database("insta").Collection("user")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		// BROKEN Finding the one user from the DB
		err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}

		// Returning a response with the user Object
		json.NewEncoder(response).Encode(user)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}

	fmt.Fprintf(response, "Endpoint Hit: Get a user by Id endpoint")
}

func postUser(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (user *User)
	   Create a user
	*/
	if request.Method == "POST" {
		response.Header().Add("content-type", "application/json")
		var user User
		json.NewDecoder(request.Body).Decode(&user)
		collection := client.Database("insta").Collection("user")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		result, _ := collection.InsertOne(ctx, user)
		json.NewEncoder(response).Encode(result)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}

	fmt.Fprintf(response, "Endpoint Hit: Create an user endpoint")
}

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
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption  string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL string             `json:"image_url,omitempty" bson:"image_url,omitempty"`
	PostedAt *time.Time         `json:"posted_at, omitempty" bson:"posted_at,omitempty"`
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
	}
	fmt.Fprintf(response, "Endpoint Hit: Get all user Posts endpoint")
}
