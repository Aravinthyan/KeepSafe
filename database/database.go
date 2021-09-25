package database

import (
	"encoding/json"
	"io/ioutil"
)

type passwordDB struct {
	WholeDB  map[string]string
	Services []string
}

func New() *passwordDB {
	db := new(passwordDB)
	db.WholeDB = nil
	db.Services = nil
	return db
}

// ReadPasswords will read the "passwords" file in pwd. The data is
// saved in a JSON format, which is extracted to be saved in a map
func (db *passwordDB) ReadPasswords() (bool, bool) {

	fileExists, validJSON := true, true

	if jsonData, err := ioutil.ReadFile("./passwords"); err == nil {
		if json.Unmarshal(jsonData, &db.WholeDB) != nil {
			validJSON = false
		}
	} else {
		fileExists = false
	}

	return fileExists, validJSON
}

// WritePasswords will write the passwords to a file called "passwords"
// in pwd.
func (db *passwordDB) WritePasswords() (bool, bool) {

	validJSON, writeError := true, true

	if jsonData, err := json.Marshal(db.WholeDB); err == nil {
		if err := ioutil.WriteFile("./passwords", jsonData, 0700); err != nil {
			writeError = false
		}
	} else {
		validJSON = false
	}

	return validJSON, writeError
}
