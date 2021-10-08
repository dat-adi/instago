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
    response.Header().Add("content-type", "application/json")
    params := request.URL.Query().Get("id")
    fmt.Println(params)
    id, _ := primitive.ObjectIDFromHex(params)
    var user User

    collection := client.Database("insta").Collection("user")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    // BROKEN for some reason
    err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{"message":"` + err.Error() + `"}`))
        return
    }

    json.NewEncoder(response).Encode(user)
	fmt.Println("Endpoint Hit: All articles endpoint")
}

func postUser(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (user *User)
	   Create a user
	*/
    response.Header().Add("content-type", "application/json")
    var user User
    json.NewDecoder(request.Body).Decode(&user)
    collection := client.Database("insta").Collection("user")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    result, _ := collection.InsertOne(ctx, user)
    json.NewEncoder(response).Encode(result)

	fmt.Println("Endpoint Hit: Create an user endpoint")
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
	ID       int64      `json:"id"`
	Caption  string     `json:"caption"`
	ImageURL string     `json:"image_url"`
	PostedAt *time.Time `json:"posted_at"`
}

type Posts []Post

func getPost(response http.ResponseWriter, request *http.Request) {
	/*
	   return type => (post *Post)
	   parameter type => (id int64)
	   Get one post by ID
	*/
	fmt.Println("Endpoint Hit: Get Post by ID endpoint")
}

func postPost(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (post *Post)
	   Create a post
	*/
    response.Header().Add("content-type", "application/json")
    var post Post
    json.NewDecoder(request.Body).Decode(&post)
    collection := client.Database("insta").Collection("user")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    result, _ := collection.InsertOne(ctx, post)
    json.NewEncoder(response).Encode(result)

	fmt.Println("Endpoint Hit: Create a Post endpoint")
}

func getAllPosts(response http.ResponseWriter, request *http.Request) {
	/*
	   paramter type => (user *User)
	   return type => array of (post *Post)?
	   Get all the posts made by a user
	*/
	fmt.Println("Endpoint Hit: Get all user Posts endpoint")
}

