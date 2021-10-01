package database

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Aravinthyan/KeepSafe/crypto"
)

// passwordDb holds the passwords for each service.
type passwordDB struct {
	WholeDB  map[string]string
	Services []string
}

// New creates a new passwordDB.
func New() *passwordDB {
	db := new(passwordDB)
	db.WholeDB = nil
	db.Services = nil
	return db
}

// ReadPasswords will read the "passwords" file in pwd and save the data in a map.
func (db *passwordDB) ReadPasswords(password []byte) error {

	var (
		encryptedJSONData []byte
		jsonData          []byte
		err               error
	)

	if encryptedJSONData, err = ioutil.ReadFile("./passwords"); err != nil {
		db.WholeDB = make(map[string]string)
		return fmt.Errorf("no existing passwords: %s", err)
	}

	if jsonData, err = crypto.Decrypt(password, encryptedJSONData); err != nil {
		return fmt.Errorf("decryption failed: %s", err)
	}

	if err = json.Unmarshal(jsonData, &db.WholeDB); err != nil {
		return fmt.Errorf("failed to convert from JSON to map: %s", err)
	}

	// create slice which is five more than the length of the whole database
	db.Services = make([]string, len(db.WholeDB), len(db.WholeDB)+5)

	index := 0
	for service := range db.WholeDB {
		db.Services[index] = service
		index++
	}

	return nil
}

// WritePasswords will write the passwords to a file called "passwords" in pwd.
func (db *passwordDB) WritePasswords(password []byte) error {

	var (
		jsonData   []byte
		outputFile *os.File
		err        error
	)

	if jsonData, err = json.Marshal(db.WholeDB); err != nil {
		return fmt.Errorf("failed to convert from map to JSON: %s", err)
	}

	var encryptedJSONData string
	if encryptedJSONData, err = crypto.Encrypt(password, jsonData); err != nil {
		return fmt.Errorf("encryption failed: %s", err)
	}

	if outputFile, err = os.Create("./passwords"); err != nil {
		return fmt.Errorf("path error: %s", err)
	}
	defer outputFile.Close()

	if err = binary.Write(outputFile, binary.LittleEndian, []byte(encryptedJSONData)); err != nil {
		return fmt.Errorf("binary.Write failed: %s", err)
	}

	return nil
}
