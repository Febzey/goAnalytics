package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

//
// Handling public keys that get assigned to a client for their hostname (website)
//

// ! Create table for public keys, public keys will be created,
// ! when a client registers on website, and requests a script link,
// ! we automatically assign a public key linked to their main token
// ! the pub key will be appended to end of script url automatically,
// ! when a analytic event happens, we assign the hostname to the public key
// ! after pub key is assigned a hostname, never allow that pub key to be used with another hostname.
// ! a pub key can only be made under and after a client is made.
// ! when assigning a new pub key, lets look in clientCache for our token who is making it,
// ! then assign the clients ID to the pub key so its linked

// Public keys are the keys attached to hostnames and used in analytic events
// The client first makes an account and recieves a Client struct, when they request
// the script for analytics, a new public key is generated, and then asSsigned the hostname on first start

// not encrypted
// we can eventually assign some permissions. (paid service etc)
// for permissions we would check the clients actual token, to see if they have paid.
type PublicKey struct {
	// unique ID for our public key
	ID int `json:"id"`

	// which client id this public key corrosponds to. (owner)
	ClientId int `json:"client_id"`

	// the actual public token to be used in requests for hostname first used with.
	PublicKey string `json:"public_key"`

	// hostname that this public key belongs to, only assigned on first use.
	Hostname string `json:"hostname"`
}

// Creating a public key that is assigned to a client.
// ! We do not assign the hostname here, hostname will be assigned automatically
// ! once the script that has this public key is loaded for the first time.
// ! we return the public key to put in our script tag
// ! this function requires the users plaintext key

// ! we will not save to pub key cache here, we will save to cache whenever it is used.
// ! user must be logged in thru our front end panel for this to work.
// ! we store a public key under a user (client) so to save a public key we need user and pass

// * The key doesnt really need to be secure since
// * it is only used to identify which hostname is using the script
// * it cant be used by another hostname once it is assigned to one.
// * and if someone does else does get and use the key, it will not work,
// * since that key is already assigned to a hostname.

func (as *AuthService) CreateAndSavePublicKey(username, password string) (string, error) {

	// get our client
	client, err := as.GetAndVerifyClient(username, password)
	if err != nil {
		return "", err
	}

	// we should to some checks here to make sure the user is allowed to make a public key
	// we should also check if the user has a public key that is not assigned to a hostname already.
	// we dont want people to make more keys than they need or are using.

	pubKey, err := generateUniqueID()
	if err != nil {
		return "", err
	}

	stmt, err := as.db.Prepare("INSERT INTO public_keys (client_id, key) VALUES(?,?)")
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	_, err = stmt.Exec(client.ID, pubKey)
	if err != nil {
		return "", err
	}

	// append the pub key to end of analytic script
	return pubKey, nil
}

// Get public key, should be used when the pub key is being used in the clients script
// we will first check cache, then check database, if not in cache. if its in database we save to cache.
// at this point we dont know if a hostname has been asigned or not, we are simply just getting what is stored in the database.
func (as *AuthService) GetPublicKeyFromCacheOrDatabase(pubKey string) (PublicKey, error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	cachedPublicKey, ok := as.publicKeyCache[pubKey]
	if !ok {
		publicKey, err := as.getPublicKeyFromDatabase(pubKey)
		if err != nil {
			return PublicKey{}, err
		}

		as.publicKeyCache[pubKey] = publicKey

		return publicKey, nil

	}

	return cachedPublicKey, nil
}

// Getting all public keys for a client,

// Getting the public key from database
func (as *AuthService) getPublicKeyFromDatabase(pubKey string) (PublicKey, error) {

	var publicKey PublicKey

	err := as.db.QueryRow("SELECT id, client_id, key, hostname WHERE key = ?", pubKey).Scan(&publicKey)
	if err != nil {
		return PublicKey{}, err
	}

	return publicKey, nil

}

// Verifying the public key, and assigning a hostname to it if it is not already assigned.
// we will check our cache first, then check the database if not in cache, we add it to cache if we find it in database,
// the function will take in the hostname and the public key, and assign the hostname to the public key.
// we will also check if the public key is already assigned to a hostname, if it is we return an error.
func (as *AuthService) VerifyPublicKey(hostname, pubKey string) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	// get the public key
	publicKey, err := as.GetPublicKeyFromCacheOrDatabase(pubKey)
	if err != nil {
		return errors.New("public key not found")
	}

	// check if the public key already has a hostname assigned
	if publicKey.Hostname != "" {
		return errors.New("public key already has a hostname assigned")
	}

	// assign the hostname to the public key
	stmt, err := as.db.Prepare("UPDATE public_keys SET hostname = ? WHERE key = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(hostname, pubKey)
	if err != nil {
		return err
	}

	// update the public key in the cache
	publicKey.Hostname = hostname
	as.publicKeyCache[pubKey] = publicKey

	return nil
}

// Generating the public key.
func generateUniqueID() (string, error) {
	// Generate a random byte slice
	idBytes := make([]byte, 16)
	_, err := rand.Read(idBytes)
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a hexadecimal string
	id := hex.EncodeToString(idBytes)

	return id, nil
}
