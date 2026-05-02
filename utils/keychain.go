package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zalando/go-keyring"
	"go.uber.org/zap"
)

// ///// Keychain Functions ///////
func KeyChainSet(service, user, secret string) (string, error) {

	keyringErr := keyring.Set(service, user, secret)
	if keyringErr != nil {
		zap.S().Error("failed to set secret: %v", keyringErr)
		return "", keyringErr
	}

	if keyringErr == nil {
		comment := "Secret set in keyring for service: " + service + " and user: " + user
		return comment, nil
	}

	return "", nil
}

func KeyChainGet(service, user string) (string, error) {
	secret, err := keyring.Get(service, user)
	if err != nil {
		zap.S().Error("failed to get secret: %v", err)
		return err.Error(), err
	}

	return secret, nil

}

func KeyChainDelete(service, user string) error {
	err := keyring.Delete(service, user)
	if err != nil {
		return err
	}
	return nil

}

func KeyChainDeleteAllServiceSecrets(service string) error {
	err := keyring.DeleteAll(service)
	if err != nil {
		return err
	}
	return nil

}

// ////////// Keychain DB Functions ////////
func KeyChainDBFilePath(filename string) (string, error) {
	dir, err := GetDirectory()
	if err != nil {
		fmt.Println(err)
	}
	path := filepath.Join(dir, filename)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return "", err
	}
	return path, nil
}

func KeychainOpenDB(filename string) (*sql.DB, error) {
	path, err := KeyChainDBFilePath(filename)
	if err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", path)

}

// EnsureSchema creates the table if it does not exist.
func KeychainDBEnsureSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS credentials (
			service TEXT NOT NULL,
			username TEXT NOT NULL,
			PRIMARY KEY (service, username)
		);
	`)
	return err
}

// SaveCredential stores service + username.
// Duplicate entries are ignored.
func KeychainDBSaveSecret(db *sql.DB, service, username string) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO credentials (service, username) VALUES (?, ?)`,
		service,
		username,
	)
	return err
}

func KeychainDBCredentialExist(db *sql.DB, username, service string) (bool, error) {
	var exists bool

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM credentials
			WHERE service = ?
			  AND username = ?
		)
	`, service, username).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func KeychainDBListServicesUsers(db *sql.DB) ([][2]string, error) {
	rows, err := db.Query(`SELECT service, username FROM credentials`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Array of exactly 2 strings - to prevent anyone trying to access the password
	var out [][2]string
	for rows.Next() {
		var service, user string
		if err := rows.Scan(&service, &user); err != nil {
			return nil, err
		}
		out = append(out, [2]string{service, user})
	}
	return out, rows.Err()
}

func KeychainDBDeleteSecret(db *sql.DB, service, username string) error {
	_, err := db.Exec(
		`DELETE FROM credentials WHERE service = ? AND username = ?`,
		service,
		username,
	)
	if err != nil {
		return err
	}

	// 2. Ensure WAL data is flushed (safe even if WAL is not enabled)
	_, err = db.Exec(`PRAGMA wal_checkpoint(FULL);`)
	if err != nil {
		return err
	}

	// 3. VACUUM must NOT run inside a transaction
	_, err = db.Exec(`VACUUM;`)
	if err != nil {
		return err
	}

	return nil
}

func KeychainDBDeleteAllSecrets(db *sql.DB) error {
	// Explicitly delete all rows
	_, err := db.Exec(`DELETE FROM credentials;`)
	if err != nil {
		return err
	}

	// 2. Ensure WAL data is flushed (safe even if WAL is not enabled)
	_, err = db.Exec(`PRAGMA wal_checkpoint(FULL);`)
	if err != nil {
		return err
	}

	// 3. VACUUM must NOT run inside a transaction
	_, err = db.Exec(`VACUUM;`)
	if err != nil {
		return err
	}

	return nil
}

func KeyChainDeleteDBFile(filename string) error {

	err := DeleteDirectory(filename)
	if err != nil {
		return err
	}
	return nil
}
