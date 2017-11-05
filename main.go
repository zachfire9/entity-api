package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type Entity struct {
    ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Date    time.Time `json:"date"`
    Type    string `json:"type"`
    Meta    interface{} `json:"meta"`
}

type Entities []Entity

var db *mgo.Database

func init() {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    db = session.DB("go")
}

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
    router.HandleFunc("/entities", EntitiesGetAll).Methods("GET")
    router.HandleFunc("/entities/{entityId}", EntitiesGet).Methods("GET")
    router.HandleFunc("/entities", EntitiesPost).Methods("POST")
    router.HandleFunc("/entities/{entityId}", EntitiesPut).Methods("PUT")
    router.HandleFunc("/entities/{entityId}", EntitiesDelete).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "You can send a RESTful request to the /entities collection.")
}

func EntitiesGet(w http.ResponseWriter, r *http.Request) {

    params := mux.Vars(r)
    entityId := params["entityId"]

    var entity Entity
    err := db.C("entities").FindId(bson.ObjectIdHex(entityId)).One(&entity)

    if err != nil {
        panic(err)
    }

    respBody, err := json.MarshalIndent(entity, "", "  ")

    if err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write(respBody)
}

func EntitiesGetAll(w http.ResponseWriter, r *http.Request) {

    var entities []Entity
    err := db.C("entities").Find(bson.M{}).All(&entities)

    if err != nil {
        panic(err)
    }

    respBody, err := json.MarshalIndent(entities, "", "  ")

    if err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write(respBody)
}

func EntitiesPost(w http.ResponseWriter, r *http.Request) {

    var entity Entity
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&entity)

    entity.ID = bson.NewObjectId()
    entity.Date = time.Now()

    if err != nil {
        panic(err)
    }

    err = db.C("entities").Insert(entity)

    if err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Location", r.URL.Path+"/"+entity.ID.Hex())
    w.WriteHeader(http.StatusCreated)
}

func EntitiesPut(w http.ResponseWriter, r *http.Request) {

    params := mux.Vars(r)
    entityId := params["entityId"]

    var entity Entity
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&entity)

    entity.Date = time.Now()

    if err != nil {
        panic(err)
    }

    err = db.C("entities").Update(bson.M{"_id": bson.ObjectIdHex(entityId)}, &entity)

    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusNoContent)
}

func EntitiesDelete(w http.ResponseWriter, r *http.Request) {

    params := mux.Vars(r)
    entityId := params["entityId"]

    err := db.C("entities").Remove(bson.M{"_id": bson.ObjectIdHex(entityId)})

    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusNoContent)
}
