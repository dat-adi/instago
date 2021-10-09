package post

import (
	"context"
	"encoding/json"
    "log"
	"fmt"
	connect "github.com/dat-adi/instago/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// Creating a struct for the Post type
type Post struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption  string             `json:"caption" bson:"caption,omitempty"`
	ImageURL string             `json:"image_url" bson:"image_url,omitempty"`
	PostedAt *time.Time         `json:"posted_at" bson:"posted_at"`
	userId   User               `json:"user_id" bson:"user_id"`
}

// Defining an array of Users
type Posts []Post

func getPost(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (http.ResponseWriter, *http.Request)
	   Get one post by ID
	*/

	// Simple logging statements
	fmt.Println("From the GET post function => ", request.RequestURI)

	// Checking the HTTP Request Method
	if request.Method == "GET" {
		// Addition of a Header to the response
		response.Header().Add("content-type", "application/json")

		// Attempting to query the :id
		// from the request URL
		params := request.URL.Query().Get("id")
		fmt.Fprintf(response, params)

		// Assigning the parsed id to a variable
		id, _ := primitive.ObjectIDFromHex(params)
		var post Post

		// Connect to the database
		collection, err := connect.getMongoDbCollection("insta", "post")

		// Assigning errors to the object and results to the :post variable
		err = collection.FindOne(context.TODO(), Post{ID: id}).Decode(&post)

		// Error Handling
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}

		// Provides indentation for a better overview
		// of the JSON object
		output, err := json.MarshalIndent(post, "", " ")
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
	fmt.Fprintf(response, "Endpoint Hit: Get Post by ID endpoint")
}

func postPost(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (http.ResponseWriter, *http.Request)
	   Create a post
	*/

	// Checks the HTTP Request Type
	if request.Method == "POST" {
		// Addition of a response header for JSON compatibility
		response.Header().Add("Content-Type", "application/json")

        // Defining an object for the post
		var post Post

		// Decodes the data and places it into the post variable
		json.NewDecoder(request.Body).Decode(&post)

		// Connects to the database
		collection, err := connect.getMongoDbCollection("insta", "post")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": Might be a problem with the Database."}`))
			return
		}

		// Insertion of the post object into the database
		result, _ := collection.InsertOne(context.Background(), post)
		json.NewEncoder(response).Encode(result)

	} else {
		// Redirects back to the homepage in the case
		// that the Request is not of the required type.
		http.Redirect(response, request, "/", http.StatusFound)
	}

	// Logging output for server side logs.
	fmt.Fprintf(response, "Endpoint Hit: Create a Post endpoint")
}

/*
    The getPostsByUserId is a function that works
    to provide all the posts created by a user in the database.
    The concept is that each post has a user ObjectID in the
    schema.
    We attempt to query the various posts by providing a filter
    which checks for the user ObjectID and append it to the output.
*/

func getPostsByUserId(response http.ResponseWriter, request *http.Request) {
	/*
	   parameter type => (http.ResponseWriter, *http.Request)
	   Get all the posts made by a user
	*/

    // Simple print for logging on the server side
	fmt.Println("From the getAllPosts function => ", request.RequestURI)

	// Checks the HTTP Request Type
	if request.Method == "GET" {
		// Addition of a response header for JSON compatibility
		response.Header().Add("content-type", "application/json")

        // Retrieval of keys from the URL
		keys := request.URL.Query()
        
		id := keys.Get("id")        // User ID parameter
		lim := keys.Get("limit")    // Limit parameter for pagination

        // Simple array of posts
		var posts Posts

		// Connecting to the database's post collection
		collection, err := connect.getMongoDbCollection("insta", "post")

        // Creating a filter to parse through the Documents
        filter := Post{userId: primitive.ObjectIDFromHex(id)}

		// Assigning errors to the object and results to the :posts variable
        err = collection.Find(context.TODO(), filter).Limit(lim).Decode(&posts)

        // Error Handling
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return
			}
			log.Fatal(err)
		}

		// Provides indentation for a better overview
		// of the JSON object
		output, err := json.MarshalIndent(posts, "", " ")
		if err != nil {
			log.Fatal(err)
		}

		// Returns the JSON output
		json.NewEncoder(response).Encode(output)
	} else {
		// Redirects back to the homepage in the case
		// that the Request is not of the required type.
		http.Redirect(response, request, "/", http.StatusFound)
	}

	// Logging output for server side logs.
	fmt.Fprintf(response, "Endpoint Hit: Get all user Posts endpoint")
}
