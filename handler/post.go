package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dp3why/mrgo/backend"
	"github.com/dp3why/mrgo/constants"
	"github.com/dp3why/mrgo/model"
	"github.com/dp3why/mrgo/service"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
    mediaTypes = map[string]string{
        ".jpeg": "image",
        ".jpg":  "image",
        ".gif":  "image",
        ".png":  "image",
        ".mov":  "video",
        ".mp4":  "video",
        ".avi":  "video",
        ".flv":  "video",
        ".wmv":  "video",
    }
)



// 1. upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Received one post request %v.\n", r)

    token := r.Context().Value("user")
    claims := token.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"] 

    // if token is not valid
    if username == nil {
        http.Error(w, "User is not authorized", http.StatusUnauthorized)
        fmt.Printf("User is not authorized\n")
        return
    }

    p := model.Post{
        Id:     uuid.New().String(),
        User:   username.(string),
        Message: r.FormValue("message"), 
    }

    file, header, err := r.FormFile("media_file")
    if err != nil {
        // return 400 bad request
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }
    defer file.Close()

    suffix := filepath.Ext(header.Filename)

    if t, ok := mediaTypes[suffix]; ok {
        p.Type = t
    } else {
        p.Type = "unknown"
    }

    log.Default().Println("Media type: ", p.Type)
    log.Default().Println("\n====================")

    ctx := context.Background()

    // Upload the file to GCS
    medialink, err := backend.UploadFileToGCS(ctx, file, constants.GCS_BUCKET, p.Id + suffix)
    if err != nil {
        log.Fatalf("Error saving file to GCS: %v", err)
         return  
    }
    p.Url = medialink
    // save post to ES
    err = service.SavePost(&p)

    if err != nil {
        http.Error(w, "Failed to save post to Elasticsearch", http.StatusInternalServerError)
        log.Default().Printf("Failed to save post to Elasticsearch %v\n", err)
        return
    }
    fmt.Printf("Post is saved to GCS and Elasticsearch: %s\n", p.Message)
}


// 2 search
func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for search")
    w.Header().Set("Content-Type", "application/json")

    user := r.URL.Query().Get("user")
    keywords := r.URL.Query().Get("keywords")

    var posts []model.Post
    var err error
    if user != "" {
        posts, err = service.SearchPostsByUser(user)
    } else {
        posts, err = service.SearchPostsByKeywords(keywords)
    }

    if err != nil {
        http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to read post from backend %v.\n", err)
        return
    }

    js, err := json.Marshal(posts)
    if err != nil {
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
    w.Write(js)
}


// 3 delete
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one delete request")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeletePost(id, username); err != nil {
        http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete post from backend %v.\n", err)
        return
    }

    fmt.Printf("Post is deleted")
}