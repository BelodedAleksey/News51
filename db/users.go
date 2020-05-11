package db

import (
	"encoding/json"
	"fmt"
	"news/log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tidwall/buntdb"
)

//User s
type User struct {
	ChatID string //key
	Name   string
	Anon   bool
}

//Post s
type Post struct {
	MessageID int //key
	User      User
	Time      time.Time
	Message   *tgbotapi.Message
	Media     *tgbotapi.MediaGroupConfig
}

var (
	dbUsers *buntdb.DB
	//map with key = message id
	Messages     = map[int]*Post{}
	LockMessages = sync.RWMutex{}
)

//Init func
func Init() {
	//init database users
	var err error
	dbUsers, err = buntdb.Open("users.db")
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[DB Users]: %s", err))
		log.Fatal(err)
	}

	err = dbUsers.SetConfig(buntdb.Config{
		SyncPolicy: buntdb.EverySecond,
	})

	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[DB Users]: %s", err))
		log.Fatal(err)
	}

	//Init database news
	dbNews, err = buntdb.Open("news.db")
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[DB News]: %s", err))
		log.Fatal(err)
	}

	err = dbNews.SetConfig(buntdb.Config{
		SyncPolicy: buntdb.EverySecond,
	})

	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[DB News]: %s", err))
		log.Fatal(err)
	}

	//load to map
	/*err = db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			user := User{}
			err := json.Unmarshal([]byte(value), &user)
			if err != nil {
				log.Errorf("[DB]: %s", err)
			}
			Users[key] = &user
			return true
		})
		return err
	})*/
}

func getEntry(user *User) (*User, error) {
	returnentry := User{}
	err := dbUsers.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(user.ChatID)
		if err != nil {
			return err
		}
		user := User{}
		err = json.Unmarshal([]byte(val), &user)
		if err != nil {
			return err
		}
		returnentry = user
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &returnentry, nil
}

//GetOrCreateEntry f
func GetOrCreateEntry(user *User) (*User, error) {
	entry, err := getEntry(user)
	if err == buntdb.ErrNotFound {
		err = addEntry(user)
		if err != nil {
			return nil, err
		}
		entry = user
	}

	return entry, nil
}

func addEntry(user *User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = dbUsers.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(user.ChatID, string(b), nil)
		return err
	})

	return err
}

//UpdateEntry f
func UpdateEntry(user *User) error {
	entry, err := GetOrCreateEntry(user)
	if err != nil {
		return err
	}

	entry.Anon = user.Anon

	if user.Name != "" {
		entry.Name = user.Name
	}

	err = addEntry(entry)
	if err != nil {
		return err
	}

	return nil
}
