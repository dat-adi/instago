package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
    "crypto/rand"
    "crypto/sha512"
	"log"
	"net/http"
	"time"
    "encoding/json"
    "encoding/base64"
)

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage Endpoint Hit...\n")
	fmt.Fprintf(response, "Please check the other endpoints now.")
    fmt.Println(request.URL)
}

/*
   Routing for the requests is done here
*/
func handleRequests() {
	http.HandleFunc("/", homePage)
    http.HandleFunc("/users/all", getAllUsers)
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

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Users []User

func getAllUsers(response http.ResponseWriter, request *http.Request) {
    fmt.Println(request)

    // Connecting to MongoDB's user collection
    collection := client.Database("insta").Collection("user")
    //ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    fmt.Println("If you see this, it's getting past here.")

    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    cursor, err := collection.Find(ctx, bson.D{})
    if err != nil {
        response.Write([]byte(`{"message":"` + err.Error() + `"}`))
    }

    var user User
    nerr := cursor.Decode(&user)
    if nerr != nil {
        response.Write([]byte(`{"message":"` + err.Error() + `"}`))
    }

    fmt.Println(cursor)
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{"message":"` + err.Error() + `"}`))
        return
    }

    json.NewEncoder(response).Encode(cursor)
}

func getUser(response http.ResponseWriter, request *http.Request) {
	/*
	   return type => (user *User)
	   parameter type => (id int64)
	   Get one user by ID
	*/

    fmt.Println(request)
    fmt.Println(request.URL.Query())

	// Checking the HTTP Request Method
	if request.Method == "GET" {
		// Addition of a Header to the response
		response.Header().Add("content-type", "application/json")

		// Parsing the URL for parameters
		params := request.URL.Query().Get("id")
        fmt.Println(params)
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

        salt := generateRandomSalt(16)
        user.Password = hashPassword(user.Password, salt)

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
	}
    fmt.Fprintf(response, "Endpoint Hit: Get all user Posts endpoint")
}

/*
    Utility Functions
*/

// Define salt size
const saltSize = 16

// Generate 16 bytes randomly and securely using the
// crypto.rand package
func generateRandomSalt(saltSize int) []byte {
  var salt = make([]byte, saltSize)

  _, err := rand.Read(salt[:])

  if err != nil {
    panic(err)
  }

  return salt
}

// Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a base64 encoded string
func hashPassword(password string, salt []byte) string {
  // Convert password string to byte slice
  var passwordBytes = []byte(password)

  // Create sha-512 hasher
  var sha512Hasher = sha512.New()

  // Append salt to password
  passwordBytes = append(passwordBytes, salt...)

  // Write password bytes to the hasher
  sha512Hasher.Write(passwordBytes)

  // Get the SHA-512 hashed password
  var hashedPasswordBytes = sha512Hasher.Sum(nil)

  // Convert the hashed password to a base64 encoded string
  var base64EncodedPasswordHash =
      base64.URLEncoding.EncodeToString(hashedPasswordBytes)

  return base64EncodedPasswordHash
}

// Check if two passwords match
func doPasswordsMatch(hashedPassword, currPassword string,
  salt[]byte) bool {
  var currPasswordHash = hashPassword(currPassword, salt)

  return hashedPassword == currPasswordHash
}
