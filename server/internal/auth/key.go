package auth

import (
	"database/sql"
	"fmt"
	"sync"
)

type AuthService struct {
	db *sql.DB
	mu *sync.Mutex

	// caching clients so database queries are not needed each time,
	// we can also easily get the client id we need
	clientCache    map[string]Client
	publicKeyCache map[string]PublicKey
}

func NewAuthService(db *sql.DB) *AuthService {
	authService := &AuthService{
		db:             db,
		mu:             &sync.Mutex{},
		clientCache:    make(map[string]Client),
		publicKeyCache: make(map[string]PublicKey),
	}
	err := authService.EnsureTables()
	if err != nil {
		fmt.Println("Error ensuring tables: ", err)

	}
	return authService
}

func (ks *AuthService) EnsureTables() error {
	// Define the table creation queries
	tableQueries := map[string]string{
		"clients": `CREATE TABLE IF NOT EXISTS clients (
			id INT AUTO_INCREMENT PRIMARY KEY,
			client_key VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL
		)`,
		"public_keys": `CREATE TABLE IF NOT EXISTS public_keys (
			id INT AUTO_INCREMENT PRIMARY KEY,
			client_id INT,
			public_key VARCHAR(255) NOT NULL UNIQUE,
			hostname VARCHAR(255) NOT NULL UNIQUE,
			FOREIGN KEY (client_id) REFERENCES clients(id)
		)`,
	}

	// Loop over the table queries and execute them
	for tableName, query := range tableQueries {
		_, err := ks.db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create %s table: %v", tableName, err)
		}
	}

	return nil
}
