package shortener

import (
	"time"
	"fmt"
	"strings"
	"database/sql"
	"encoding/base64"
)

var ConfigGl Config

// Contains fields that will be saved to database
type URL struct {
	Alias     string
	URL       string
	Timestamp time.Time
}

func generateShortURL(longURL string) (string, error) {

	u := URL{
		"",
		longURL,
		time.Now().UTC(),
	}

	// Open a transaction that will be passed to underline functions to prevent race condition while writing to database
	tx, err := DB.Begin()
	if err != nil {
		err = fmt.Errorf("Could not start sql transaction. %v\n", err)
		return "", err
	}

	err = u.createAlias(tx)
	if err != nil {
		err = fmt.Errorf("Could not get alias. %v\n", err)
		return "", err
	}

	u.normalizeURL()

	err = u.saveToDB(tx)
	if err != nil {
		err = fmt.Errorf("Could not save to database. %v\n", err)
		return "", err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("Could not commit transaction. %v\n", err)
		return "", err
	}

	shortURL := fmt.Sprintf("%v/url/%v", ConfigGl.URLHost, u.Alias)

	return shortURL, nil
}

func (u *URL) createAlias(tx *sql.Tx) error {
	var alias string

	// Get a hash for shortening <longURL + timestamp>
	hash := u.Hash()

	// Simple solution to work around hash collisions
	// Since we are not using the entire hash we can use different sub parts of it in case of a collision
	for i := 0; i < len(hash); i++ {
		end := i + ConfigGl.URLAliasLength

		if end > len(hash) {
			break
		}

		alias = hash[i:end]

		isUsed, err := checkAlias(alias, tx)
		if err != nil {
			err = fmt.Errorf("Could not check if alias is already in use. %v\n", err)
			return err
		}

		if !isUsed {
			u.Alias = alias
			return nil
		}
	}

	err := fmt.Errorf("A unique alias could not be created")
	return err
}

func (u *URL) Hash() string {
	// Using timestamp as a salt to create a unique hash from the same URL
	s := u.URL + u.Timestamp.String()

	hash := base64.StdEncoding.EncodeToString([]byte(s))

	return hash
}

func (u *URL) normalizeURL() {
	s := u.URL

	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		u.URL = "http://" + s
	}
}
