package resources

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

// shaders_block_fragmentshader reads file data from disk. It returns an error on failure.
func shaders_block_fragmentshader() ([]byte, error) {
	return bindata_read(
		"c:\\goworkspace-win\\wnd\\src\\wnd\\resources\\shaders\\block.fragmentshader",
		"shaders/block.fragmentshader",
	)
}

// shaders_block_vertexshader reads file data from disk. It returns an error on failure.
func shaders_block_vertexshader() ([]byte, error) {
	return bindata_read(
		"c:\\goworkspace-win\\wnd\\src\\wnd\\resources\\shaders\\block.vertexshader",
		"shaders/block.vertexshader",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"shaders/block.fragmentshader": shaders_block_fragmentshader,
	"shaders/block.vertexshader": shaders_block_vertexshader,
}
