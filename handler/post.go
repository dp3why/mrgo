package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

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

    // 1. process the request
    // form data -> go struct
    p := model.Post{
        Id:      uuid.New().String(),
        User:    r.FormValue("user"),
        Message: r.FormValue("message"),
    }
 
    file, header, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        log.Default().Printf("Media file is not available %v\n", err)
        return
    }
 
    //p.Type
    suffix := filepath.Ext(header.Filename)
    if t, ok := mediaTypes[suffix]; ok {
        p.Type = t
    } else {
        p.Type = "unknown"
    }
 
    // 2. call service to handle request
    err = service.SavePost(&p, file)
    if err != nil {
        http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
        log.Default().Printf("Failed to save post to backend %v\n", err)
        return
    }
 
    // 3. construct response
    log.Default().Println("Post is saved successfully.")
 }
 

// 2 search
func searchHandler(w http.ResponseWriter, r *http.Request) {
    log.Default().Println("Received one request for search")
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
        log.Default().Printf("Failed to read post from backend %v.\n", err)
        return
    }

    js, err := json.Marshal(posts)
    if err != nil {
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        log.Default().Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
    w.Write(js)
}


// 3 delete
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    log.Default().Println("Received one delete request")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeletePost(id, username); err != nil {
        http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
        log.Default().Printf("Failed to delete post from backend %v.\n", err)
        return
    }

    log.Default().Printf("Post is deleted")
}