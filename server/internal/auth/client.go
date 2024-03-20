package auth

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

//
// Handling clients. (people who signed up for our service)
//

// ! our website will have a front end panel to hit endpoint to create a new client.
// ! after they create account they will be able to get their public key by hitting a buttn on the front end.

// private key that will be linked with the client who signed up
type Client struct {
	// Id generated for our client by database.
	ID int `json:"id"`

	// Client Main Key
	ClientKey string `json:"client_key"`

	// email for the client
	Email string `json:"email"`

	//Encrypted AES Password
	Password string `json:"password"`
}

// Creating and Saving API key to database,
// returns the plaintext key to give back to token owner.
// Key linked to a hostname.
// ! we should add optional passowrd if the use our auth system.
func (as *AuthService) CreateAndSaveClient(email, password string) error {
	stmt, err := as.db.Prepare("INSERT INTO clients (client_key, email, password) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing SQL statement: ", err)
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	fmt.Println("Creating and saving client")

	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return err
	}

	fmt.Println(" encryptedPassword: ", encryptedPassword)

	clientKey, err := generateUniqueID()
	if err != nil {
		return err
	}

	fmt.Println(" clientKey: ", clientKey)

	// Execute the SQL statement
	_, err = stmt.Exec(clientKey, email, encryptedPassword)
	if err != nil {
		// Check for duplicate key error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			// Duplicate key error occurred
			return fmt.Errorf("exists")
		}
		return fmt.Errorf("internal error: %v", err)
	}

	fmt.Println("Client created and saved successfully")
	return nil
}

// Getting an api key from database or cache, by given plainTextKey.
// we check the cache for the token before looking into the database.
func (as *AuthService) GetAndVerifyClient(email, plainTextPass string) (Client, error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	var client Client

	cachedClient, ok := as.clientCache[email]
	if !ok {
		clientFromDb, err := as.GetClientFromDatabase(email)
		if err != nil {
			return Client{}, err
		}
		as.clientCache[clientFromDb.Email] = clientFromDb

		client = clientFromDb
	} else {
		client = cachedClient
	}

	match, err := comparePasswordsAES(plainTextPass, client.Password)
	if err != nil {
		return Client{}, err
	}

	fmt.Println("match: ", match)

	if match {
		return client, nil
	} else {
		return Client{}, fmt.Errorf("badpass")
	}
}

// Get API KEY from database
func (as *AuthService) GetClientFromDatabase(email string) (Client, error) {

	var client Client

	err := as.db.QueryRow("SELECT * FROM clients WHERE email = ?", email).Scan(
		&client.ID,
		&client.ClientKey,
		&client.Email,
		&client.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Client{}, fmt.Errorf("not found")
		}
		return Client{}, err
	}

	return client, nil
}
