package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
  //  "io/ioutil"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    )
    
var client *mongo.Client

type User struct {
    ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    UserId         string         `json:"user_id,omitempty" bson:"user_id,omitempty"`
    Name             string      `json:"name,omitempty" bson:"name,omitempty"`
    PhoneNumber      string      `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
    EmailAddress     string      `json:"emailaddress,omitempty" bson:"emailaddress,omitempty"`
    TimeStamp        string      `json:"timeStamp" bson: "timeStamp"`
    DateOfBirth      string      `json:"dob,omitempty" bson: "dob,omitempty"`
}

type Contact struct{
    CID        primitive.ObjectID   `json:"c_id,omitempty" bson:"c_id,omitempty"`
    UserOneId       string          `json:"user_id_one,omitempty" bson:"user_id_one,omitempty"`
    UserTwoId       string          `json:"user_id_two,omitempty" bson:"user_id_two,omitempty"`
    TimeOfContact   string          `json:"contact_time,omitempty" bson:"contact_time,omitempty"`
}

//functio to add contact between two user
func AddContact(response http.ResponseWriter, request *http.Request){
     response.Header().Set("content-type", "application/json")
     var contact Contact
     _=json.NewDecoder(request.Body).Decode(&contact)
     collection:=client.Database("ContactTracingApi").Collection("Contacts")
     ctx,_:=context.WithTimeout(context.Background(), 5*time.Second)
     result,_:=collection.InsertOne(ctx,contact)
     fmt.Println(result)
     json.NewEncoder(response).Encode(contact)
}


// function to get the complete list of user
func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
    var users []User
    collection := client.Database("ContactTracingApi").Collection("users")
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
        return
    }
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var user User
        cursor.Decode(&user)
        users = append(users, user)
    }
    if err := cursor.Err(); err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
        return
    }
    json.NewEncoder(response).Encode(users)
}

//function to create a new user
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
    var user User
    _ = json.NewDecoder(request.Body).Decode(&user)
    collection := client.Database("ContactTracingApi").Collection("users")
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
     user.TimeStamp = time.Now().String()
     result,_ := collection.InsertOne(ctx, user)
     fmt.Println(result)
    json.NewEncoder(response).Encode(user)
}


func main() {

    fmt.Println("Starting the application...")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    clientOptions := options.Client().ApplyURI("mongodb+srv://kanak:appointy@cluster0.ivuws.mongodb.net/Data")
    client, _ = mongo.Connect(ctx, clientOptions)
    //to create a new user
    http.HandleFunc("/users", CreateUserEndpoint)  

    //to get the list of all the users
    http.HandleFunc("/user", GetUsersEndpoint)
    
    //to add contact between to user
    http.HandleFunc("/contact",AddContact)

    log.Fatal(http.ListenAndServe(":12345",nil))
}

