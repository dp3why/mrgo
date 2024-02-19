package handler

import (
	"net/http"

	jwtMiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitRouter() http.Handler {
	jwtMiddleware := jwtMiddleware.New(jwtMiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router:= mux.NewRouter()
	
	router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
	router.Handle("/search", http.HandlerFunc(searchHandler)).Methods("GET")
	router.Handle("/post/{id}", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")
	
	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signHandler)).Methods("POST")

	originsOK := handlers.AllowedOrigins([]string{"*"})
	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	return handlers.CORS(originsOK, headersOK, methodsOK)(router)
}
