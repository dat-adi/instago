package model

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

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
