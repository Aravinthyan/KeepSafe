package database

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Aravinthyan/KeepSafe/crypto"
)

// PasswordDB holds the passwords for each service.
type PasswordDB struct {
	wholeDB  map[string]string
	Services []string
}

const PasswordFile = "./passwords"

// New creates a new PasswordDB.
func New() *PasswordDB {
	db := new(PasswordDB)
	db.wholeDB = nil
	db.Services = nil
	return db
}

// ReadPasswords will read the "passwords" file in pwd and save the data in a map.
func (db *PasswordDB) ReadPasswords(password []byte) error {

	var (
		encryptedJSONData []byte
		jsonData          []byte
		err               error
	)

	if encryptedJSONData, err = ioutil.ReadFile(PasswordFile); err != nil {
		db.wholeDB = make(map[string]string)
		return fmt.Errorf("no existing passwords: %s", err)
	}

	if jsonData, err = crypto.Decrypt(password, encryptedJSONData); err != nil {
		return fmt.Errorf("decryption failed: %s", err)
	}

	if err = json.Unmarshal(jsonData, &db.wholeDB); err != nil {
		return fmt.Errorf("failed to convert from JSON to map: %s", err)
	}

	// create slice which is five more than the length of the whole database
	db.Services = make([]string, len(db.wholeDB), len(db.wholeDB)+5)

	index := 0
	for service := range db.wholeDB {
		db.Services[index] = service
		index++
	}

	return nil
}

// WritePasswords will write the passwords to a file called "passwords" in pwd.
func (db *PasswordDB) WritePasswords(password []byte) error {

	var (
		jsonData   []byte
		outputFile *os.File
		err        error
	)

	if db.wholeDB == nil {
		return errors.New("nothing to write")
	}

	if jsonData, err = json.Marshal(db.wholeDB); err != nil {
		return fmt.Errorf("failed to convert from map to JSON: %s", err)
	}

	var encryptedJSONData string
	if encryptedJSONData, err = crypto.Encrypt(password, jsonData); err != nil {
		return fmt.Errorf("encryption failed: %s", err)
	}

	if outputFile, err = os.Create(PasswordFile); err != nil {
		return fmt.Errorf("path error: %s", err)
	}
	defer outputFile.Close()

	if err = binary.Write(outputFile, binary.LittleEndian, []byte(encryptedJSONData)); err != nil {
		return fmt.Errorf("binary.Write failed: %s", err)
	}

	return nil
}
