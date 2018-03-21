package shortener

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

var Db *sql.DB

const (
	DbDriver     = "postgres"
	DbTimeFormat = "2006-01-02 15:04:05"
)

func CreateDBConnection(c *Config) error {
	host := c.DbConfig.DbHost
	port := c.DbConfig.DbPort
	user := c.DbConfig.DbUsername
	pass := c.DbConfig.DbPassword
	db := c.DbConfig.DbName

	config := "host=" + host + " port=" + port + " dbname=" + db + " user=" + user + " password=" + pass

	var err error
	Db, err = sql.Open(DbDriver, config)
	if err != nil {
		err = fmt.Errorf("Could not open sql connection. %v\n", err)
		return err
	}

	return nil
}

// Saves fields of the Url struct to database
func (u Url) saveToDB(tx *sql.Tx) error {

	table := ConfigGl.DbConfig.Table

	query := "INSERT INTO " + table + " (alias, url, timestamp) VALUES ($1, $2, $3);"
	params := []interface{}{u.Alias, u.Url, u.Timestamp.Format(DbTimeFormat)}

	// Execute the insert statement
	_, err := tx.Exec(query, params...)
	if err != nil {
		err = fmt.Errorf("Could not execute sql statement '%v'. %v\n", query, err)
		return err
	}

	return nil
}

func getLongUrl(alias string) (string, error) {

	var longUrl string

	table := ConfigGl.DbConfig.Table

	query := "SELECT url FROM " + table + " WHERE alias = $1;"

	// Execute the select sql
	row := Db.QueryRow(query, alias)

	err := row.Scan(&longUrl)
	if err != nil {
		err = fmt.Errorf("Could not run query '%v'. %v\n", query, err)
		return "", err
	}

	return longUrl, nil
}

func checkAlias(alias string, tx *sql.Tx) (bool, error) {

	table := ConfigGl.DbConfig.Table

	query := "SELECT * FROM " + table + " WHERE alias = $1;"

	// Execute the select statement
	_, err := tx.Query(query, alias)
	if err != nil {
		// If alias is not present in DB
		if err == sql.ErrNoRows {
			return true, nil
		}
		// If Query returned any other error
		return false, err
	}

	// If alias already in use
	return false, nil
}
