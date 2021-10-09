package model

import (
    "testing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUsers(t *testing.T) {
	firstUser := User{ primitive.NewObjectID(), "Dat Adi", "adihata@gmail.com", "something"}
	secondUser := User{ primitive.NewObjectID(), "Aaron", "abcd@gmail.com", "nothing"}
	thirdUser := User{ primitive.NewObjectID(), "Swagg", "aww@gmail.com", "nottingdale"}
}
