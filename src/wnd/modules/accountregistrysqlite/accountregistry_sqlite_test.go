package accountregistrysqlite

import (
	"wnd/utils/test"
	
	"os"
    "testing"
)

var ar *accountRegistrySQLite

func TestCreate(t *testing.T) {
	ar = new(accountRegistrySQLite)
	ar.GameDir = test.NewGameDir()
	
	os.Remove(ar.getDBLocation())
	
	_,err1 := ar.Login("testuser", "testpassword", true)
	
	if err1 != nil {
		t.Error(err1)
		t.FailNow()
	}
	
	err2 := ar.Logout("testuser")
	if err2 != nil {
		t.Error(err2)
		t.FailNow()
	}
	
	_,err3 := ar.Login("testuser", "testpassword", false)
	if err3 != nil {
		t.Error(err3)
		t.FailNow()
	}
}

