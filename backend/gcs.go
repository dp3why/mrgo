package backend

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/alexsasharegan/dotenv"
	"google.golang.org/api/option"
)


func UploadFileToGCS(ctx context.Context, file io.Reader, bucketName, objectName string) (string, error) {
	_ = dotenv.Load()

	encryptedCredentialsPath := os.Getenv("ENCRYPTED_CREDENTIALS_PATH")
	privateKeyBase64 := os.Getenv("PRIVATE_KEY_BASE64")
	if encryptedCredentialsPath == "" || privateKeyBase64 == "" {
		 return "", fmt.Errorf("ENCRYPTED_CREDENTIALS_PATH or PRIVATE_KEY_BASE64 environment variable is not set")
	}

	// Using the DecryptCredentials function to get decrypted credentials
	decryptedCredentials :=  DecryptCredentials(encryptedCredentialsPath, privateKeyBase64)

	// Creating a new storage client using the decrypted credentials
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsJSON(decryptedCredentials))
	if err != nil {
		 log.Default().Panicf("failed to create storage client with decrypted credentials: %v", err)
		 return "", err
	}
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

