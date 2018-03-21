package shortener

import (
	"time"
	"encoding/hex"
	"crypto/md5"
	"fmt"
	"strings"
	"database/sql"
)

var ConfigGl Config

// Contains fields that will be saved to database
type Url struct {
	Alias     string
	Url       string
	Timestamp time.Time
}

func generateShortUrl(longUrl string) (string, error) {

	hash := getMd5Hash(longUrl, time.Now())

	// Open a transaction that will be passed to underline functions to prevent race condition while writing to database
	tx, err := Db.Begin()
	if err != nil {
		err = fmt.Errorf("Could not start sql transaction. %v\n", err)
		return "", err
	}

	alias, err := getAlias(hash, tx)
	if err != nil {
		err = fmt.Errorf("Could not get alias. %v\n", err)
		return "", err
	}

	normalizeUrl(&longUrl)

	url := Url{
		Alias:     alias,
		Url:       longUrl,
		Timestamp: time.Now().UTC(),
	}

	err = url.saveToDB(tx)
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

	shortUrl := ConfigGl.UrlHost + "/url/" + alias

	return shortUrl, nil
}

func getMd5Hash(s string, t time.Time) string {
	s = s + t.String()

	hasher := md5.New()
	hasher.Write([]byte(s))
	hashed := hex.EncodeToString(hasher.Sum(nil))

	return hashed
}

func getAlias(s string, tx *sql.Tx) (string, error) {
	var alias string

	// Simple solution to work around hash collisions
	// Since we are not using the entire hash we can use different sub parts of it in case of a collision
	for i := 0; i < len(s); i++ {
		end := i + ConfigGl.UrlAliasLength

		if end > len(s) {
			break
		}

		alias = s[i:end]

		isUsed, err := checkAlias(alias, tx)
		if err != nil {
			err = fmt.Errorf("Could not check if alias is already in use. %v\n", err)
			return "", err
		}

		if !isUsed {
			return alias, nil
		}
	}

	err := fmt.Errorf("A unique alias could not be created")
	return "", err
}

func normalizeUrl(s *string) {
	if !strings.HasPrefix(*s, "http://") && !strings.HasPrefix(*s, "https://") {
		*s = "http://" + *s
	}
}
