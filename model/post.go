package model

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	connect "github.com/dat-adi/instago/database/connect"
)

type Post struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Caption  string             `json:"caption" bson:"caption,omitempty"`
	ImageURL string             `json:"image_url" bson:"image_url,omitempty"`
	PostedAt *time.Time         `json:"posted_at" bson:"posted_at"`
    userId   User               `json:"user_id" bson:"user_id"`
}

type Posts []Post

func getPost(response http.ResponseWriter, request *http.Request) {
	/*
	   return type => (post *Post)
	   parameter type => (id int64)
	   Get one post by ID
	*/
    fmt.Println("From the GET post function => ", request.RequestURI)

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
    fmt.Println("From the POST post function => ", request.RequestURI)

	if request.Method == "POST" {
		response.Header().Add("content-type", "application/json")
		var post Post
		json.NewDecoder(request.Body).Decode(&post)

        collection, err := connect.getMongoDbCollection("insta", "post")
        if err != nil{
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": Might be a problem with the Database."}`))
            return
        }

		result, _ := collection.InsertOne(context.Background(), post)
		json.NewEncoder(response).Encode(result)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}

	fmt.Fprintf(response, "Endpoint Hit: Create a Post endpoint")
}

func getPostsByUserId(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (user *User)
	   return type => array of (post *Post)?
	   Get all the posts made by a user
	*/

    fmt.Println("From the getAllPosts function => ", request.RequestURI)

	if request.Method == "GET" {
        response.Header().Add("content-type", "application/json")
        keys := request.URL.Query()
        fmt.Println("Are we reaching here?")
        fmt.Printf("Here are the keys : %s\n", keys)
        id := keys.Get("id")
        fmt.Println(id)
        //lim := keys.Get("limit")

        //var posts []Post
        collection, err := connect.getMongoDbCollection("insta", "post")
        cursor, err := collection.Find(context.TODO(), bson.M{"userId": id})
        if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": Might be a problem with the Database."}`))
            return
        }

        fmt.Println(cursor)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
	fmt.Fprintf(response, "Endpoint Hit: Get all user Posts endpoint")
}
