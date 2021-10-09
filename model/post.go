package model

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
