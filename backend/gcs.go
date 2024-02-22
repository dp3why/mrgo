package backend

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/dp3why/mrgo/constants"
	"google.golang.org/api/option"
)

var (
   GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
   client *storage.Client
   bucket string
}

func InitGCSBackend() {

	ctx := context.Background()

	// Getting the encrypted credentials path and private key base64 from the environment variables
	privateKeyBase64 := os.Getenv("PRIVATE_KEY_BASE64")
	if constants.EncryptedCredentialsPath == "" || privateKeyBase64 == "" {
		log.Default().Panicf("ENCRYPTED_CREDENTIALS_PATH or PRIVATE_KEY_BASE64 is not set")
		return 
	}

	// Using the DecryptCredentials function to get decrypted credentials
	decryptedCredentials :=  DecryptCredentials(constants.EncryptedCredentialsPath, privateKeyBase64)

	// Creating a new storage client using the decrypted credentials
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(decryptedCredentials))
	if err != nil {
		log.Default().Panicf("failed to create storage client with decrypted credentials: %v", err)
		return
	}
	log.Default().Println("storage client connected.")

	GCSBackend = &GoogleCloudStorageBackend{
		client: client,
		bucket: constants.GCS_BUCKET,
	}
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
   ctx := context.Background()
   object := backend.client.Bucket(backend.bucket).Object(objectName)
   wc := object.NewWriter(ctx)
   if _, err := io.Copy(wc, r); err != nil {
       return "", err
   }

   if err := wc.Close(); err != nil {
       return "", err
   }

   if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
       return "", err
   }

   attrs, err := object.Attrs(ctx)
   if err != nil {
       return "", err
   }

   log.Default().Printf("File is saved to GCS: %s\n", attrs.MediaLink)
   return attrs.MediaLink, nil
}

func DeleteFromGS(objectName string) error {
   ctx := context.Background()
   object := GCSBackend.client.Bucket(GCSBackend.bucket).Object(objectName)
   if err := object.Delete(ctx); err != nil {
	   return err
   }
   log.Default().Println("File is deleted from GCS")
   return nil
}