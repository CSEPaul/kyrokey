package internal

import (
	"fmt"
	"kyrokey/utils"

	"go.uber.org/zap"
)

func KcDeleteSecret(service string, user string) string {
	err := utils.KeyChainDelete(service, user)
	if err != nil {
		zap.S().Error("failed to delete secret: %v", err)
	}

	fmt.Println("Secret deleted from keyring related to:--- service:", service, "and user:", user)

	//delete the service + user entry from sqlite db
	db, err := utils.KeychainOpenDB(utils.KeyChainDB)
	if err != nil {
		zap.S().Error(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			zap.S().Error(err)
		}
	}()

	err = utils.KeychainDBDeleteSecret(db, service, user)
	if err != nil {
		zap.S().Error(err)
	}
	fmt.Println("All users related to service:", service, "deleted from kc_cli tracking db.")
	return "Secret deleted from kc_cli tracking db."
}

func KcDeleteDB(confirm string) string {
	dir, err := utils.KeyChainDBFilePath(utils.KeyChainDB)
	if err != nil {
		zap.S().Error("Error getting db file path:", err.Error())
	}

	comment := "Secret DB Deleted"
	comment2 := "Write the `Confirm` flag"
	switch confirm {
	case "Confirm", "C":
		err := utils.KeyChainDeleteDBFile(dir)
		if err != nil {
			zap.S().Error("Error deleting db file:", err.Error())
		}
		// check if db is present - keychain_entries.db
		exists := utils.FileExists(utils.KeyChainDB)
		if exists {
			statement := "DB Not Deleted"
			return statement
		}
		return comment

	default:
		println("You must use the Confirm flag to delete the db file for security.")
		return comment2
	}

}

func KcGet(service string, user string) string {
	var secret string

	secret, err := utils.KeyChainGet(service, user)
	if err != nil {
		zap.S().Error("failed to get secret: %v", err)
	}
	return secret
}

func KcList() [][2]string {
	db, err := utils.KeychainOpenDB(utils.KeyChainDB)
	if err != nil {
		zap.S().Error(err)
	}
	defer db.Close()

	entries, err := utils.KeychainDBListServicesUsers(db)
	if err != nil {
		zap.S().Error("failed to list kc_cli entries: %v", err)
	}

	return entries
}

func KcSet(service string, user string, secret string) error {

	/*
		Check the keychaindb to see if that service and user already exist
		If they do the return error - of already used credientials
	*/
	db, err := utils.KeychainOpenDB(utils.KeyChainDB)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			zap.S().Error(err)
		}
	}()

	if err := utils.KeychainDBEnsureSchema(db); err != nil {
		zap.S().Error(err)
		return err
	}
	exists, err := utils.KeychainDBCredentialExist(db, user, service)
	if err != nil {
		zap.S().Error(err)
		return nil
	}
	if exists {
		zap.S().Info("Using existing credentials for service " + service)
		// stop this operation, return control to GUI/CLI
		return fmt.Errorf("service/user already exists")
	}

	// Set the Secret in the real Keychain
	comment, err := utils.KeyChainSet(service, user, secret)
	if err != nil {
		zap.S().Error("failed to set secret: %v", err)
		// exit the app to stop the adding to the database
		return err
	}
	fmt.Println(comment)

	//write the service name and user to a sqlite db to track kc_cli entries
	if err := utils.KeychainDBSaveSecret(db, service, user); err != nil {
		return err
	}

	return nil

}
