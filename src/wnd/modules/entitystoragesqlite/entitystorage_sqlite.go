package entitystoragesqlite

import (
	"os"
	"wnd/api"
	"wnd/modules"
	"wnd/utils/log"

	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const (
	_DB_NAME = "entities.wndb"

	_CREATE_TABLE_QUERY  = "create table entities (id text not null primary key, xposition float, yposition float, zposition float, xlooksat float, ylooksat float, zlooksat float, control text)"
	_CHECK_TABLE_EXISTS  = "SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'entities'"
	_ADD_ENTITY_QUERY    = "insert into entities(id, xposition, yposition, zposition, xlooksat, ylooksat, zlooksat, control) values(?, ?, ?, ?, ?, ?, ?, ?)"
	_REMOVE_ENTITY_QUERY = "delete from entities where id = ?"

	_SELECT_QUERY                   = "select id, xposition, yposition, zposition, xlooksat, ylooksat, zlooksat, control from entities where %s = ?"
	_SELECT_BY_POSITION_RANGE_QUERY = "select id, xposition, yposition, zposition, xlooksat, ylooksat, zlooksat, control from entities where xposition >= ? and xposition <= ? and yposition >= ? and yposition <= ? and zposition >= ? and zposition <= ?"
)

type entityStorageSQLite struct {
	sync.RWMutex

	GameDir  modules.GameDir `inject:""`
	Universe *api.Universe   `inject:""`
}

func (this *entityStorageSQLite) ID() string {
	return "entityStorageSQLite"
}

func (this *entityStorageSQLite) Add(e *api.Entity) error {
	log.Tracef("entityStorageSQLite.Add: %v", e)

	this.Lock()
	defer this.Unlock()

	if err := this.checkFileExists(); err != nil {
		return err
	}

	log.Debugf("db created")

	db, err := this.open()
	if err != nil {
		return err
	}

	defer db.Close() //may error

	_, err3 := db.Exec(_ADD_ENTITY_QUERY, e.ID, e.Position.X, e.Position.Y, e.Position.Z, e.LooksAt.X, e.LooksAt.Y, e.LooksAt.Z, e.Control)
	return err3
}

func (this *entityStorageSQLite) Remove(id string) error {
	log.Tracef("entityStorageSQLite.Remove: %s", id)

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

	_, err2 := db.Exec(_REMOVE_ENTITY_QUERY, id)
	return err2
}

func (this *entityStorageSQLite) get(paramName string, paramValue interface{}) ([]*api.Entity, error) {
	if err := this.checkFileExists(); err != nil {
		return nil, err
	}

	db, err := this.open()
	if err != nil {
		return nil, err
	}

	defer db.Close() //may error

	res, err2 := db.Query(fmt.Sprintf(_SELECT_QUERY, paramName), paramValue)
	if err2 != nil {
		return nil, err2
	}

	entities := make([]*api.Entity, 0)
	for res.Next() {
		entity := new(api.Entity)

		if e := res.Scan(&entity.ID, &entity.Position.X, &entity.Position.Y, &entity.Position.Z, &entity.LooksAt.X, &entity.LooksAt.Y, &entity.LooksAt.Z, &entity.Control); e != nil {
			return nil, e
		} else {
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

func (this *entityStorageSQLite) GetByID(id string) (*api.Entity, error) {
	log.Tracef("entityStorageSQLite.GetByID: %s", id)

	this.Lock()
	defer this.Unlock()

	res, err := this.get("id", id)

	if err != nil {
		return nil, err
	}

	if res == nil || len(res) < 1 {
		return nil, errors.New("No entity found with ID# " + id)
	}

	return res[0], nil
}
func (this *entityStorageSQLite) GetByCoords(coords api.WorldCoords, radius uint32) ([]*api.Entity, error) {
	log.Tracef("entityStorageSQLite.GetBtCoords: %v (%d)", coords.Pack(), radius)

	this.Lock()
	defer this.Unlock()

	db, err := this.open()
	if err != nil {
		return nil, err
	}

	defer db.Close() //may error

	res, err2 := db.Query(_SELECT_BY_POSITION_RANGE_QUERY, coords.X-int32(radius), coords.X+int32(radius), coords.Y-int32(radius), coords.Y+int32(radius), coords.Z-int32(radius), coords.Z+int32(radius))
	if err2 != nil {
		return nil, err2
	}

	entities := make([]*api.Entity, 0)
	for res.Next() {
		entity := new(api.Entity)

		if e := res.Scan(&entity.ID, &entity.Position.X, &entity.Position.Y, &entity.Position.Z, &entity.LooksAt.X, &entity.LooksAt.Y, &entity.LooksAt.Z, &entity.Control); e != nil {
			return nil, e
		} else {
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

func (this *entityStorageSQLite) GetByControl(control string) ([]*api.Entity, error) {
	log.Tracef("entityStorageSQLite.GetByControl: %s", control)

	this.Lock()
	defer this.Unlock()

	return this.get("control", control)
}

func (this *entityStorageSQLite) getDBLocation() string {
	return this.GameDir.UniverseSavesDir(this.Universe.Name) + "/" + _DB_NAME
}

func (this *entityStorageSQLite) checkFileExists() error {
	log.Tracef("entityStorageSQLite.checkFileExists")

	if _, err := os.Stat(this.getDBLocation()); os.IsNotExist(err) {
		log.Infof("entityStorageSQLite: creating new entities file")

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

func (this *entityStorageSQLite) open() (*sql.DB, error) {
	return sql.Open("sqlite3", this.getDBLocation())
}
