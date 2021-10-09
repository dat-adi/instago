package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dat-adi/instago/database/connect"
	passec "github.com/dat-adi/instago/utils/passec"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

// Creating a struct for the User type
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

// Defining an array of Users
type Users []User

// Defining the client object for the database
var client *mongo.Client

/*
The getUser is a function that is created to find a particular
user in the database, according to the ID.
*/
func getUser(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (http.ResponseWriter, *http.Request)
	   Get one user by ID
	*/

	// Checking the HTTP Request Method
	if request.Method == "GET" {
		// Addition of a Header to the response
		response.Header().Set("Content-Type", "application/json")

        // Attempting to query the :id
        // from the request URL
		params := request.URL.Query().Get("id")
		fmt.Fprintf(response, params)

        // Assigning the parsed id to a variable
		id, _ := primitive.ObjectIDFromHex(params)
        var user User

		// Connecting to MongoDB's user collection
		collection, err := connect.getMongoDbCollection("insta", "user")

		// Assigning errors to the object and results to the :user variable
		err = collection.FindOne(context.TODO(), User{ID: id}).Decode(&user)

		// Error Handling
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return
			}
			log.Fatal(err)
		}

		// Provides indentation for a better overview
		// of the JSON object
		output, err := json.MarshalIndent(user, "", " ")
		if err != nil {
			log.Fatal(err)
		}

		// Returns the JSON output
		json.NewEncoder(response).Encode(output)
	} else {
		// Redirects back to the homepage in the case
		// that the Request is not of the required type.
		http.Redirect(response, request, "/users/all", http.StatusFound)
	}

	// Logging output for server side logs.
	fmt.Fprintf(response, "Endpoint Hit: Get a user by Id endpoint")
}

func postUser(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (http.ResponseWriter, *http.Request)
	   Create a users
	*/

	// Checks the HTTP Request Type
	if request.Method == "POST" {
		// Addition of a response header for JSON compatibility
		response.Header().Set("Content-Type", "application/json")

		// Defining an object for the user
		var user User

		// Decodes the data and places it into the user variable
		json.NewDecoder(request.Body).Decode(&user)

		// Connects to the database
		collection, err := connect.getMongoDbCollection("insta", "post")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": Might be a problem with the Database."}`))
			return
		}

		// Generates a salt and encrypts the password
		salt := passec.generateRandomSalt(16)
		user.Password = passec.hashPassword(user.Password, salt)

		// Insertion of the user object into the database
		result, _ := collection.InsertOne(context.Background(), user)
		json.NewEncoder(response).Encode(result)

	} else {
		// Redirects back to the homepage in the case
		// that the Request is not of the required type.
		http.Redirect(response, request, "/", http.StatusFound)
	}

	// Logging output for server side logs.
	fmt.Fprintf(response, "Endpoint Hit: Create an user endpoint")
}

/*
    The getAllUsers is a function that works
    to provide all the users in the database.
    * Not a part of the given tasks.
*/

func getAllUsers(response http.ResponseWriter, request *http.Request) {
	// Creating a simple user object
	var users Users

	// Connecting to MongoDB's user collection
    collection, err := connect.getMongoDbCollection("insta", "user")

	// Retrieving the collection of users
	err = collection.Find(context.TODO(), users).Decode(&users)

    // Error Handling
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Providing the list of users in the response
	json.NewEncoder(response).Encode(users)
}
