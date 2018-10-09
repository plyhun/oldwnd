package entitystoragesqlite

import (
	"wnd/api"
	"wnd/utils/test"

	"os"
	"testing"
)

var ar *entityStorageSQLite

func TestCreate(t *testing.T) {
	ar = new(entityStorageSQLite)
	ar.GameDir = test.NewGameDir()
	ar.Universe = new(api.Universe)
	ar.Universe.Name = "Uni"

	os.Remove(ar.getDBLocation())

	e := &api.Entity{
		ID:      "entity1",
		Control: "player1",
		Position: api.Coords{
			X: -12.2,
			Y: -13.3,
			Z: -14.4,
		},
		LooksAt: api.Coords{
			X: 13.4,
			Y: 14.5,
			Z: 15.6,
		},
	}

	err1 := ar.Add(e)

	if err1 != nil {
		t.Error(err1)
		t.FailNow()
	}

	e1, err2 := ar.GetByID("entity1")
	if err2 != nil {
		t.Error(err2)
		t.FailNow()
	} else {
		t.Logf("GetByID: %v", e1)
	}

	e2, err3 := ar.GetByControl("player1")
	if err3 != nil {
		t.Error(err3)
		t.FailNow()
	} else {
		t.Logf("GetByControl: %v", e2[0])
	}

	position := api.WorldCoords{
		X: -16,
		Y: -16,
		Z: -16,
	}

	e3, err4 := ar.GetByCoords(position, 10)
	if err4 != nil {
		t.Error(err4)
		t.FailNow()
	} else {
		t.Logf("GetByCoords: %v", e3[0])
	}

	position = api.WorldCoords{
		X: -56,
		Y: -56,
		Z: -56,
	}

	e4, err5 := ar.GetByCoords(position, 10)
	if err5 != nil {
		t.Error(err5)
		t.FailNow()
	} else {
		t.Logf("GetByCoords: %v", e4)
	}
}
