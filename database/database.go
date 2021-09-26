package database

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Aravinthyan/KeepSafe/crypto"
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

	return nil
}

// WritePasswords will write the passwords to a file called "passwords"
// in pwd.
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
