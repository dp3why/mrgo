package backend

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

// Gen generates a new RSA key pair and uses the public key to encrypt a new AES key.
func Gen(jsonFilePath string, encryptedFilePath string) {
    // Generate RSA Keys
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        panic(err)
    }

    // Convert the private key to PEM format
    privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: privateKeyBytes,
    })

    // Encode the PEM-encoded private key to a Base64 string
    privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)
	// Write the Base64 encoded private key to a .env file
	envContent := fmt.Sprintf("PRIVATE_KEY_BASE64=%s\n", privateKeyBase64)
	if err := os.WriteFile("./.env", []byte(envContent), 0600); err != nil {
		panic(err)
	}
    fmt.Println("Base64 Encoded Private Key saved to .env file")

 

    // Read the JSON file
    plaintext, err := os.ReadFile(jsonFilePath)
    if err != nil {
        panic(err)
    }

    // Generate a new AES key for encrypting the data
    aesKey := make([]byte, 32) // AES-256
    if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
        panic(err)
    }

    block, err := aes.NewCipher(aesKey)
    if err != nil {
        panic(err)
    }

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    // Encrypt the AES key with the RSA public key
    encryptedAesKey, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, aesKey)
    if err != nil {
        panic(err)
    }

    // Save the encrypted AES key and the encrypted data to the file
    // First part of the file is the encrypted AES key, then the IV + ciphertext.
    finalData := append(encryptedAesKey, ciphertext...)
    if err := os.WriteFile(encryptedFilePath, finalData, 0600); err != nil {
        panic(err)
    }

    fmt.Println("Encrypted data written to:", encryptedFilePath)
}

// DecryptCredentials decrypts the encrypted credentials file using the RSA private key.
func DecryptCredentials(filePath, privateKeyBase64 string) []byte {
    privateKeyPEM, err := base64.StdEncoding.DecodeString(privateKeyBase64)
    if err != nil {
        panic(err)
    }

    pemBlock, _ := pem.Decode(privateKeyPEM)
    if pemBlock == nil {
        panic("failed to parse PEM block containing the private key")
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
    if err != nil {
        panic(err)
    }

    encryptedData, err := os.ReadFile(filePath)
    if err != nil {
        panic(err)
    }

    encryptedAesKey := encryptedData[:256]
    aesKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedAesKey)
    if err != nil {
        panic(err)
    }

    aesCipher, err := aes.NewCipher(aesKey)
    if err != nil {
        panic(err)
    }

    iv := encryptedData[256 : 256+aes.BlockSize]
    ciphertext := encryptedData[256+aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(aesCipher, iv)
    plaintext := make([]byte, len(ciphertext))
    stream.XORKeyStream(plaintext, ciphertext)

    return plaintext // Correctly returning decrypted content
}




