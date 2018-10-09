package accountregistrysqlite

import (
	"wnd/api"
	"wnd/modules"
	"wnd/utils/log"

	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const (
	_DB_NAME = "accounts.wndb"

	_CREATE_TABLE_QUERY     = "create table accounts (id text not null primary key, universes text, pwhash text, loggedin boolean)"
	_MODIFY_LOGIN_QUERY     = "update accounts set loggedin = ? where id = ?"
	_MODIFY_UNIVERSES_QUERY = "update accounts set universes = ? where id = ?"
	_MODIFY_PWHASH_QUERY    = "update accounts set pwhash = ? where id = ?"
	_DELETE_QUERY           = "delete from accounts where id = ?"
	_CHECK_TABLE_EXISTS     = "SELECT name FROM sqlite_master WHERE type='table' AND name='accounts'"
	_ADD_ACCOUNT_QUERY      = "insert into accounts(id, pwhash, universes, loggedin) values(?, ?, ?, ?)"

	_SELECT_PW_QUERY = "select pwhash, universes, loggedin from accounts where id=?"
	_SELECT_QUERY    = "select universes, loggedin from accounts where id = ?"
	_SELECT_ALL      = "select id, universes, loggedin from accounts"
)

func New() modules.AccountRegistry {
	return &accountRegistrySQLite{}
}

type accountRegistrySQLite struct {
	sync.RWMutex

	GameDir modules.GameDir `inject:""`
}

func (this *accountRegistrySQLite) ID() string {
	return "accountRegistrySQLite"
}

func (this *accountRegistrySQLite) Login(id, password string, createIfAbsent bool) (*api.Account, error) {
	log.Tracef("%s, create %v", id, createIfAbsent)

	this.Lock()
	defer this.Unlock()

	return this.loginInternal(id, password, createIfAbsent)
}

func (this *accountRegistrySQLite) Logout(id string) error {
	log.Tracef("%s", id)

	this.Lock()
	defer this.Unlock()

	if err := this.checkFileExists(); err != nil {
		return err
	}

	return this.modifyLoginState(id, false)
}

func (this *accountRegistrySQLite) Delete(id string) error {
	log.Tracef("%s", id)

	this.Lock()
	defer this.Unlock()

	if err := this.checkFileExists(); err != nil {
		return err
	}

	db, err := this.open()
	if err != nil {
		return err
	}

	defer db.Close() //may error

	_, err2 := db.Exec(_DELETE_QUERY, id)
	return err2
}

func (this *accountRegistrySQLite) Get(id string) (*api.Account, error) {
	log.Tracef("%s", id)

	this.RLock()
	defer this.RUnlock()

	if err := this.checkFileExists(); err != nil {
		return nil, err
	}

	db, err := this.open()
	if err != nil {
		return nil, err
	}

	defer db.Close() //may error

	res, err2 := db.Query(_SELECT_QUERY, id)
	if err2 != nil {
		return nil, err2
	}

	universes := ""
	count := 0
	loggedin := false

	for res.Next() {
		count++

		if count > 1 {
			continue
		}

		if e := res.Scan(&universes, &loggedin); e != nil {
			return nil, e
		}
	}

	if count < 1 {
		return nil, errors.New("Account " + id + "not found")
	} else {
		return &api.Account{
			ID:        id,
			Universes: unwrapUniverses(universes),
		}, nil
	}
}

func (this *accountRegistrySQLite) Update(a *api.Account) error {
	log.Tracef("%v", a)

	this.Lock()
	defer this.Unlock()

	if err := this.checkFileExists(); err != nil {
		return err
	}

	db, err := this.open()
	if err != nil {
		return err
	}

	defer db.Close() //may error

	_, err2 := db.Query(_MODIFY_UNIVERSES_QUERY, wrapUniverses(a.Universes), a.ID)
	return err2
}

func (this *accountRegistrySQLite) UpdatePassword(id, oldpw, newpw string) error {
	log.Tracef("%v", id)

	this.Lock()
	defer this.Unlock()

	if _, err := this.loginInternal(id, oldpw, false); err != nil {
		return err
	}

	return this.modify(_MODIFY_PWHASH_QUERY, hashPw(newpw), id)
}

func (this *accountRegistrySQLite) modifyLoginState(id string, state bool) error {
	log.Tracef("%s => %v", id, state)

	return this.modify(_MODIFY_LOGIN_QUERY, state, id)
}

func (this *accountRegistrySQLite) modify(query string, params ...interface{}) error {
	log.Tracef("%s & %v", query, params)

	if err := this.checkFileExists(); err != nil {
		return err
	}

	db, err := this.open()
	if err != nil {
		return err
	}

	defer db.Close() //may error

	_, err2 := db.Exec(query, params...)
	return err2
}

func (this *accountRegistrySQLite) checkFileExists() error {
	log.Tracef("")

	if _, err := os.Stat(this.getDBLocation()); os.IsNotExist(err) {
		log.Infof("creating new accounts file")

		f, err2 := os.Create(this.getDBLocation())
		if err2 != nil {
			return err2
		}

		f.Close()
		log.Debugf("created %v", f.Name())
	}

	db, err3 := this.open()
	if err3 != nil {
		return err3
	}

	log.Debugf("opened %v", db)

	defer db.Close() //may also error

	res, err4 := db.Query(_CHECK_TABLE_EXISTS)
	if err4 != nil {
		return err4
	}

	count := 0
	for res.Next() {
		count++
	}

	log.Debugf("tables: %v", count)

	if count < 1 {
		_, err5 := db.Exec(_CREATE_TABLE_QUERY)
		return err5
	}

	return nil
}

func (this *accountRegistrySQLite) open() (*sql.DB, error) {
	return sql.Open("sqlite3", this.getDBLocation())
}

func (this *accountRegistrySQLite) getDBLocation() string {
	return this.GameDir.StorageDir() + "/" + _DB_NAME
}

func (this *accountRegistrySQLite) loginInternal(id, password string, createIfAbsent bool) (*api.Account, error) {
	log.Tracef("%s, create %v", id, createIfAbsent)

	if err := this.checkFileExists(); err != nil {
		return nil, err
	}

	db, err := this.open()
	if err != nil {
		return nil, err
	}

	defer db.Close() //may error

	res, err2 := db.Query(_SELECT_PW_QUERY, id)
	if err2 != nil {
		return nil, err2
	}

	var universes, pwhash string
	count := 0
	loggedin := false

	for res.Next() {
		count++

		if count > 1 {
			continue
		}

		if e := res.Scan(&pwhash, &universes, &loggedin); e != nil {
			return nil, e
		}
	}

	if loggedin {
		return nil, errors.New("Account " + id + " already logged in")
	}

	if count > 0 {
		if hashPw(password) == pwhash {
			err3 := this.modifyLoginState(id, true)
			if err3 != nil {
				return nil, err3
			}

			return &api.Account{
				ID:        id,
				Universes: unwrapUniverses(universes),
			}, nil
		} else {
			return nil, errors.New("Wrong password")
		}
	} else if createIfAbsent {
		_, err3 := db.Exec(_ADD_ACCOUNT_QUERY, id, hashPw(password), "", true)

		var account *api.Account
		if err3 == nil {
			account = &api.Account{
				ID:        id,
				Universes: []string{},
			}
		}

		return account, err3
	} else {
		return nil, errors.New("No account " + id + " found")
	}
}

func hashPw(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func wrapUniverses(universes []string) string {
	log.Tracef("%v", universes)
	
	if universes == nil {
		return ""
	}

	res := ""

	for i, u := range universes {
		if u != "" {
			res += u

			if i < (len(universes) - 1) {
				res += "/"
			}
		}
	}

	return res
}

func unwrapUniverses(universes string) []string {
	log.Tracef("%v", universes)
	
	return strings.Split(universes, "/")
}
