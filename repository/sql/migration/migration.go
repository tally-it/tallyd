// Code generated by go-bindata.
// sources:
// repository/sql/migration/resources/001_initial_db.sql
// DO NOT EDIT!

package migration

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _resources001_initial_dbSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x57\x5d\x73\xe2\x36\x14\x7d\xf7\xaf\xb8\x6f\x0b\xd3\x6c\x67\xc9\x6c\x66\xda\xe9\xec\x83\xb1\x95\xd4\xb3\x20\x52\x23\x77\x9a\x27\xad\xc0\x26\xeb\x01\x6c\x6a\xcb\x99\xe6\xdf\x77\xe4\x4f\xc9\x16\x41\x26\xed\xf2\x04\x92\x8e\x7c\x3f\xce\x39\x5c\x7f\xfc\x08\x3f\x1d\xe3\xe7\x8c\xf1\x08\x82\x93\xe5\xf8\xc8\x26\x08\x88\x3d\x5f\x20\xd8\x32\x1e\x3d\xa7\x59\x1c\xe5\xd6\xc4\x82\xe6\xe7\x2b\x8d\x43\xf0\x30\x01\x3b\x20\x2b\xea\x61\xc7\x47\x4b\x84\x89\x05\x00\xf0\xe8\x7b\x4b\xdb\x7f\x82\xaf\xe8\xe9\xc6\x02\x48\xd8\x31\x82\xfa\x23\x10\x9a\x0f\x5e\x11\xc0\xc1\x62\x21\x8e\xbf\xc4\x79\xbc\x39\x54\x08\xe2\xe1\x27\x0f\x93\xc9\x6c\x0a\x2e\xba\xb7\x83\x05\x81\x0f\xb3\x0f\xca\x71\xb6\xe5\xf1\x4b\x7d\xbf\xee\xf8\x27\xf5\x78\x9c\xd3\x2c\x4d\xb9\xc9\xed\xd6\xd4\x02\x40\xf8\xc1\xc3\x08\xbe\x80\x97\x24\xa9\x3b\xb7\x00\x9c\xdf\x6d\x7f\x8d\x08\x7c\x81\x82\xef\x7e\x39\x6e\x3e\xff\x66\x69\x0b\xf6\x4a\x4f\x2c\x8b\x12\x4e\x8f\xec\xd4\x56\x2e\x7d\xae\x2a\x27\x95\x43\x2d\xa0\x88\xb1\xc6\xf5\x2b\x2d\xa7\x21\x95\x18\x26\xd2\xc5\x37\x1a\xf0\xf4\xda\x44\x4e\x59\x1a\x16\x5b\xe9\xae\x26\x93\x66\x43\x24\xa2\x4f\xc1\x38\xf6\xee\xaa\x1b\x19\x34\x15\x07\x9d\x15\x5e\x13\xdf\x16\x78\x5d\x28\x34\xde\xec\xf6\x74\x66\x01\xdc\xaf\x7c\xe4\x3d\xe0\xae\x1a\xcd\x2d\xe0\xa3\x7b\xe4\x23\xec\xa0\xb5\xc4\x63\xf5\x4c\x49\xd9\x15\x86\xe0\xd1\x15\xa9\x3b\xf6\xda\xb1\x5d\xd4\xac\xba\x68\x81\xa4\xd5\xd1\x95\xf4\xb0\x8b\xfe\xba\x10\xfe\x0a\x6b\x0f\xa8\x61\x9e\x69\x4e\x3e\x68\x88\xae\x25\x3a\x59\xae\xbf\x06\x3d\x26\x5e\xf8\xc8\x3d\x94\x35\x0d\xf0\xa7\xed\x8b\x12\x4c\x6e\xef\xee\xa6\x06\xe0\x07\xe2\xe1\x6e\xa7\x44\xce\x3e\x9f\x03\x96\xe0\x1a\x78\xca\xe2\x6d\xf7\x58\x17\x39\xde\xd2\x5e\x4c\x66\x77\x37\x70\xab\x87\x2b\x46\x11\x86\x51\x48\x19\xaf\xc1\x36\x41\xc4\x5b\xa2\x56\xf9\x4e\xe0\xfb\x08\x13\x2a\x16\xd7\xc4\x5e\x3e\x2a\xe0\x30\x3a\x44\xbc\x85\xb7\xe0\xcb\x21\xc7\x39\x95\x0c\x6d\xee\x91\xf6\x81\x1b\xe1\x35\x6f\x87\xfc\x77\xc1\x12\x1e\xf3\x57\xf3\x7c\x7b\x40\x5a\x24\x31\x37\x6a\xcf\x55\x7e\xa7\x90\x3b\xa7\x25\xa3\x92\x30\xfa\x47\xe1\x74\x0e\x93\x8a\x6a\x82\xc2\xf6\x82\x20\xff\x2d\x7b\x01\xb0\x5d\xd7\x4c\xf9\xb7\x96\xa2\xfb\x4e\x03\x8a\xec\xbb\x28\xa4\x03\x96\x5e\xf1\x03\xbd\xf7\x45\x97\xf3\x74\xbb\x2f\x15\x57\x7e\x3b\x63\x80\x6f\xab\xcd\x40\x69\x0a\x09\x8a\x3c\xca\x2a\xa8\x11\x72\x40\x9d\xb1\xcf\x63\x61\x23\x94\x4e\x0a\x66\x22\x91\xba\x56\x95\x47\x61\xc6\x6e\xdf\x77\xea\x9a\x16\xfa\x6e\xd5\x9b\x3f\xc2\x9d\xcf\x05\xbb\xc2\xd5\x96\x4c\x60\x0d\x50\xf6\xf1\xfa\x7c\xdd\xb3\x81\x69\xf3\x8c\x25\xb9\x98\x58\xd2\xa4\x32\x6e\x69\xc1\x7c\xa0\xea\x18\x01\xc6\x84\xaa\x5b\xa4\x9a\xfe\x28\xe8\x0b\x3b\x14\x92\xe7\x1b\xb8\x91\x42\x0e\xce\x9e\xe5\x2d\x03\x4f\x3a\xe3\xdc\xa3\x69\x59\x9c\x42\xd6\x79\xb7\x04\xbf\x9c\xb3\xc4\xe8\x6f\x72\xef\x68\xc3\x96\x9f\xab\x82\xd2\xdd\xfe\xdb\xff\xc6\x6e\xc1\xdd\xf7\xf8\x73\x4d\x97\x8a\xa0\x72\x16\x3a\x9e\x56\x10\x93\x64\x07\x97\x0d\x44\x52\x71\xbe\x7c\x06\x2b\xf8\xf7\x8a\xf1\xed\xcf\x2b\xf9\xae\x50\x56\x6e\xf3\x31\xe2\xdf\xd3\x50\x47\x2f\xe5\xed\x42\x66\xb1\x78\x01\x98\x2f\x56\xf3\xb6\xe7\xd7\x96\x78\xb7\xa7\x5d\x9a\xe5\xd7\xaa\x40\xdd\xe2\x79\x4f\x10\xeb\x52\x69\xc6\x8c\x70\xd2\x14\xd6\x24\x3c\xfb\x75\x66\x32\x82\x45\x47\x16\x1f\xcc\x91\x35\x6a\x9b\x45\x8d\x92\x46\xce\x4f\x92\x06\xc7\x0d\x4f\x9b\x43\xba\xdd\x47\x61\x6f\x72\xfa\x74\x69\x72\x8a\x73\xca\xc2\x63\x9c\x0c\x66\xae\x8b\x48\x49\xf2\x65\x67\xa8\x28\x32\x2d\x9a\xb9\x26\xc0\xde\x1f\x01\x82\x89\x58\x1d\xf3\x62\x25\x0f\x3e\xd5\x14\x31\x98\x74\x94\x7f\x13\xc5\x4a\x1a\xea\xc8\x5e\x52\x06\xd7\x6d\x29\x13\x4d\xeb\x18\x3a\x1b\xe9\x05\xa3\xfc\x25\x0d\x62\xda\xed\xa9\x62\x05\x15\x59\xdf\x15\x9c\x6e\xdc\x3a\x13\x9b\x64\x1d\xba\xc8\xfa\x92\xfb\xcf\xc3\x6a\xa7\xc0\x7f\x03\x00\x00\xff\xff\x56\x14\x7e\x56\x20\x11\x00\x00")

func resources001_initial_dbSqlBytes() ([]byte, error) {
	return bindataRead(
		_resources001_initial_dbSql,
		"resources/001_initial_db.sql",
	)
}

func resources001_initial_dbSql() (*asset, error) {
	bytes, err := resources001_initial_dbSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/001_initial_db.sql", size: 4384, mode: os.FileMode(436), modTime: time.Unix(1518103360, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
var _bindata = map[string]func() (*asset, error){
	"resources/001_initial_db.sql": resources001_initial_dbSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"resources": &bintree{nil, map[string]*bintree{
		"001_initial_db.sql": &bintree{resources001_initial_dbSql, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

