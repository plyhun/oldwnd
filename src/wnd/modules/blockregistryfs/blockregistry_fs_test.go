package blockregistryfs

import (
	"wnd/utils/test"
	
	"testing"
)

var br *blocksRegistryFs

func init() {
	br = new(blocksRegistryFs)
	br.GameDir = test.NewGameDir()
}

func TestBlockRegistryInit(t *testing.T) {
	e := br.Init()
	if e != nil {
		t.Errorf("TestBlockRegistryInit failed: %v", e)
	}
}

func TestBlockRegistryStatistics(t *testing.T) {
	n,s := br.Statistics()
	t.Logf("%d blocks: %s", n, s)
}

func TestBlockRegistryAll(t *testing.T) {
	ret := br.All()
	t.Logf("Blocks: %#v", ret)
}