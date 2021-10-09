package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dat-adi/instago/database/connect"
	passec "github.com/dat-adi/instago/utils/passec"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Users []User

var client *mongo.Client

func getAllUsers(response http.ResponseWriter, request *http.Request) {
	fmt.Println(request)

	var user User

	// Connecting to MongoDB's user collection
	collection := client.Database("insta").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(ctx, user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(response).Encode(cursor)
}

func getUser(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (http.ResponseWriter, *http.Request)
	   Get one user by ID
	*/

	// Checking the HTTP Request Method
	if request.Method == "GET" {
		// Addition of a Header to the response
		response.Header().Set("Content-Type", "application/json")

		// Connecting to MongoDB's user collection
		collection, err := connect.getMongoDbCollection("insta", "user")

		var result bson.M
		err = collection.FindOne(context.TODO(), bson.D{{"name", "Datta Adithya"}}).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return
			}
			log.Fatal(err)
		}

		output, err := json.MarshalIndent(result, "", " ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", output)

		json.NewEncoder(response).Encode(output)
	} else {
		http.Redirect(response, request, "/users/all", http.StatusFound)
	}

	fmt.Fprintf(response, "Endpoint Hit: Get a user by Id endpoint")
}

func postUser(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (user *User)
	   Create a user
	*/

	fmt.Println(request)

	if request.Method == "POST" {
		response.Header().Add("content-type", "application/json")
		var user User
		json.NewDecoder(request.Body).Decode(&user)
		collection := client.Database("insta").Collection("user")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		salt := passec.generateRandomSalt(16)
		user.Password = passec.hashPassword(user.Password, salt)

		result, _ := collection.InsertOne(ctx, user)
		json.NewEncoder(response).Encode(result)
		/*
		   The function to check for whether the passwords match,
		   should be one that is done during the authentication process,
		   and as such, has not been implemented here.
		   However, can be done using the :doPasswordsMatch: function
		   in the utilities.
		*/
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}

	fmt.Fprintf(response, "Endpoint Hit: Create an user endpoint")
}
