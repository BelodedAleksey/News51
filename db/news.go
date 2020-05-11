package db

import (
	"github.com/tidwall/buntdb"
)

var (
	dbNews *buntdb.DB
)

//GetLastNew f
func GetLastNew(key string) (string, error) {
	returnentry := ""
	err := dbNews.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		returnentry = val
		return nil
	})

	if err != nil {
		return "", err
	}
	return returnentry, nil
}

//UpdateLastNew f
func UpdateLastNew(key, value string) error {
	err := dbNews.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})

	return err
}
