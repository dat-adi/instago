package main

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage Endpoint Hit...\n")
	fmt.Fprintf(response, "Please check the other endpoints now.")
    fmt.Println("From the homepage => ", request.URL)
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
	http.HandleFunc("/posts/users/{id}", getPostsByUserId)
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

	var user User

	// Connecting to MongoDB's user collection
	collection := client.Database("insta").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(ctx, user)
	if err != nil {
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
	}

	err = cursor.Decode(&user)
	if err != nil {
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

	// Checking the HTTP Request Method
	if request.Method == "GET" {
		// Addition of a Header to the response
		response.Header().Set("Content-Type", "application/json")

		// Connecting to MongoDB's user collection
        collection, err := getMongoDbCollection("insta", "user")
        if err != nil{
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": Might be a problem with the Database."}`))
            return
        }

        var filter bson.M = bson.M{}

        if request.URL.Query().Get("id") != "" {
            id := request.URL.Query().Get("id")
            objID, _ := primitive.ObjectIDFromHex(id)
            filter = bson.M{"_id": objID}
        }

        fmt.Println(request)

        var results []bson.M
        cur, err := collection.Find(context.Background(), filter)
        defer cur.Close(context.Background())

        if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
            return
        }

        cur.All(context.Background(), &results)
        if results == nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message":"` + err.Error() + `"}`))
            return
        }

		json.NewEncoder(response).Encode(results)
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

        collection, err := getMongoDbCollection("insta", "post")
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
        collection, err := getMongoDbCollection("insta", "post")
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
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}

// Check if two passwords match
func doPasswordsMatch(hashedPassword, currPassword string, salt []byte) bool {
	var currPasswordHash = hashPassword(currPassword, salt)

	return hashedPassword == currPasswordHash
}

func GetMongoDbCollection() (*mongo.Client, error) {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.Background(), readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }
    
    return client, nil
}

func getMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error){
    client, err := GetMongoDbCollection()
    if err != nil {
        return nil, err
    }

    collection := client.Database(DbName).Collection(CollectionName)

    return collection, nil
}
