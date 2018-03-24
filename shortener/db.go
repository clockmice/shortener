package shortener

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	DBDriver     = "postgres"
	DBTimeFormat = "2006-01-02 15:04:05"
)

func CreateDBConnection(c *Config) error {
	host := c.DBConfig.DBHost
	port := c.DBConfig.DBPort
	user := c.DBConfig.DBUsername
	pass := c.DBConfig.DBPassword
	db := c.DBConfig.DBName

	config := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v", host, port, db, user, pass)

	var err error
	DB, err = sql.Open(DBDriver, config)
	if err != nil {
		err = fmt.Errorf("Could not open sql connection. %v\n", err)
		return err
	}

	return nil
}

// Saves fields of the URL struct to database
func (u URL) saveToDB(tx *sql.Tx) error {

	query := fmt.Sprintf("INSERT INTO %v (alias, url, timestamp) VALUES ($1, $2, $3);", ConfigGl.DBConfig.Table)
	params := []interface{}{
		u.Alias,
		u.URL,
		u.Timestamp.Format(DBTimeFormat),
	}

	// Execute the insert statement
	_, err := tx.Exec(query, params...)
	if err != nil {
		err = fmt.Errorf("Could not execute sql statement '%v' %v. %v\n", query, params, err)
		return err
	}

	return nil
}

func getLongURL(alias string) (string, error) {

	var longURL string

	table := ConfigGl.DBConfig.Table

	query := fmt.Sprintf("SELECT url FROM %v WHERE alias = $1;", table)

	// Execute the select sql
	row := DB.QueryRow(query, alias)

	err := row.Scan(&longURL)
	if err != nil {
		err = fmt.Errorf("Could not run query '%v'. %v\n", query, err)
		return "", err
	}

	return longURL, nil
}

func checkAlias(alias string, tx *sql.Tx) (bool, error) {

	table := ConfigGl.DBConfig.Table

	query := fmt.Sprintf("SELECT * FROM %v WHERE alias = $1;", table)

	// Execute the select statement
	rows, err := tx.Query(query, alias)
	defer rows.Close()
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
