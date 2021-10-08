package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Aravinthyan/KeepSafe/crypto"
)

// PasswordDB holds the passwords for each service.
type PasswordDB struct {
	wholeDB  map[string]string
	Services []string
}

// PasswordFile contains the name of the file that contains the passwords.
const PasswordFile = "./appdata"

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
		jsonData []byte
		err      error
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

	if err = ioutil.WriteFile(PasswordFile, []byte(encryptedJSONData), 0400); err != nil {
		return fmt.Errorf("write file failed: %s", err)
	}

	return nil
}

func (db *PasswordDB) CreateEmptyDB() {
	db.wholeDB = make(map[string]string)
}

func (db *PasswordDB) Insert(service, password string) error {
	if _, exist := db.wholeDB[service]; exist {
		return errors.New("service does exist")
	}

	db.wholeDB[service] = password
	db.Services = append(db.Services, service)
	return nil
}

func (db *PasswordDB) Remove(service string) error {
	if _, exist := db.wholeDB[service]; exist {
		delete(db.wholeDB, service)

		for index, value := range db.Services {
			if value == service {
				servicesLength := len(db.Services) - 1
				db.Services[index] = db.Services[servicesLength]
				db.Services[servicesLength] = ""
				db.Services = db.Services[:servicesLength]
			}
		}
		return nil
	}

	return errors.New("service does not exist")
}

func (db *PasswordDB) Password(service string) string {
	if password, exist := db.wholeDB[service]; exist {
		return password
	}
	return ""
}
