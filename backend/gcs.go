package backend

import (
	"context"
	"io"
	"log"

	"cloud.google.com/go/storage"
	"github.com/dp3why/mrgo/constants"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)


func UploadFileToGCS(ctx context.Context, file io.Reader, bucketName, objectName string) (string, error) {
	 
	// Construct the JWT config
	config := &jwt.Config{
		Email:      constants.CLIENT_EMAIL,
		PrivateKey: []byte(constants.PRIVATE_KEY),
		Scopes:     []string{storage.ScopeFullControl},
		TokenURL:   google.JWTTokenURL,
	}

	// Create an HTTP client using the JWT config
	httpClient := config.Client(ctx)

	// Use the HTTP client to create a storage client
	storageClient, err := storage.NewClient(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		log.Default().Printf("failed to create storage client: %v.\n", err)
		return "", err
	}
    //
	log.Default().Println("storage client connected.")

	defer storageClient.Close()

	// Get the bucket and object
	bucket := storageClient.Bucket(bucketName)
	object := bucket.Object(objectName)

	// Create a new writer for the object
	wc := object.NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		log.Default().Printf("io.Copy: %v.\n", err)
		return "", err
	}
	if err := wc.Close(); err != nil {
		log.Default().Printf("Writer.Close: %v.\n", err)
		return "", err
	}

	// set ACL to read for all public users
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
        return "", err
    }

	attrs, err := object.Attrs(ctx)
    if err != nil {
		log.Default().Printf("Attrs: %v.\n", err)
        return "", err
    }

    log.Default().Printf("File is saved to GCS: %s\n", attrs.MediaLink)
    return attrs.MediaLink, nil
}

