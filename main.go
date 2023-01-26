package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func generateFromPassword(password string, p *argonParams, username string) (hash []byte, err error) {
	// user username as salt which is kinda randomish..... but atleast different for eacb user?
	salt := []byte(username)
	if err != nil {
		return nil, err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash = argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return hash, nil
}

func main() {
	// USAGE:
	// ./gocrypt e ./PathToFolder Username
	// ./gocrypt d ./PathToFile Username

	// Check and handle all usage errors
	// Take password and use PBKDF to turn into an AES key
	// if decrypt decrypt file and unzip it rinse and repeat for all subdirs

	// password := os.Args[4]

	// if len(password) < 10 {
	// 	RaiseError("Password must be greater than 10 characters", true)
	// }
	mode, PathToFileOrFolder, username := UsageParamsCheck(os.Args)
	// GENERATE PRIVATE KEY WITH ARGON2
	hash, err := generateFromPassword(
		"password123",
		&argonParams{
			memory:      64 * 1024,
			iterations:  3,
			parallelism: 2,
			saltLength:  16,
			keyLength:   32,
		},
		username,
	)

	if err != nil {
		RaiseError(fmt.Sprintf("%d", err), true)
	}

	if mode == "e" {
		// ZIP FILE/FOLDER
		fileNameArr := strings.Split(PathToFileOrFolder, "/")
		fileName := fileNameArr[len(fileNameArr)-1]
		if !strings.Contains(fileName, ".") {
			RecursiveZip(PathToFileOrFolder, PathToFileOrFolder+".zip")
			EncryptFile(fileName+".zip", hash, PathToFileOrFolder+".zip")
		} else {
			EncryptFile(fileName, hash, PathToFileOrFolder)
		}
	} else if mode == "d" {

		fileNameArr := strings.Split(PathToFileOrFolder, "/")
		encrytedFile := fileNameArr[len(fileNameArr)-1]
		fileNameExtArr := strings.Split(encrytedFile, ".")
		if strings.Contains(encrytedFile, "zip") {
			f, _ := os.Open(PathToFileOrFolder)
			data, err := ioutil.ReadAll(f)
			if err != nil {
				RaiseError("Error: Canot read encrypted file data", true)
			}

			decryptedData, err2 := Decrypt(hash, data)
			if err2 != nil {
				RaiseError("Error:  Canot decrypt data", true)
			}

			err3 := os.WriteFile(fileNameExtArr[0]+"."+fileNameExtArr[1], decryptedData, 0644)
			if err3 != nil {
				RaiseError("Error: Canot write zip", true)
			}

			Unzip(fileNameExtArr[0]+"."+fileNameExtArr[1], fileNameExtArr[0])
		} else {
			f, _ := os.Open(PathToFileOrFolder)
			data, err := ioutil.ReadAll(f)
			if err != nil {
				RaiseError("Error: Canot read encrypted file data", true)
			}

			decryptedData, err2 := Decrypt(hash, data)
			if err2 != nil {
				RaiseError("Error:  Canot decrypt data", true)
			}

			err3 := os.WriteFile(fileNameExtArr[0]+"."+fileNameExtArr[1], decryptedData, 0644)
			if err3 != nil {
				RaiseError("Error: Canot write zip", true)
			}
		}
	}
}

// new plan
// only zip if its a folder
func EncryptFile(fileName string, hash []byte, PathToFileOrFolder string) {
	f, _ := os.Open(PathToFileOrFolder)
	data, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		RaiseError(fmt.Sprintf("%d", err2), true)
	}

	encryptedData, err3 := Encrypt(hash, data)
	if err3 != nil {
		fmt.Println("1")
		fmt.Println(err3)
		RaiseError(fmt.Sprintf("%d", err3), true)
	}

	// WRITE NEW FILE WITH TIMESTAMP AS NAME I THINKS?
	err4 := os.WriteFile(fileName+".crypt", encryptedData, 0644)
	if err4 != nil {
		fmt.Println("1")
		fmt.Println(err4)
		RaiseError(fmt.Sprintf("%d", err4), true)
	}

	// DELETE OLD FILE
	e := os.RemoveAll(PathToFileOrFolder)
	if e != nil {
		fmt.Println("1")
		fmt.Println(e)
		RaiseError(fmt.Sprintf("%d", e), true)
	}
}
