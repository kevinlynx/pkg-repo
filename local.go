package pkgrepo

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LocalGetter struct {
	root string
}

func NewLocalGetter(url0 string) (*LocalGetter, error) {
	root := strings.TrimPrefix(url0, "file://")
	return &LocalGetter{root: root}, nil
}

func (self *LocalGetter) List(name string, ver string) ([]*Package, error) {
	files, err := filepath.Glob(path.Join(self.root, name+"*"))
	if err != nil {
		return nil, err
	}
	pkgs := make([]*Package, 0, len(files))
	for _, fname := range files {
		_, filename := path.Split(fname)
		checksum, _ := md5sumFile(fname) // better to lazy md5
		pkg := &Package{Name: filename, URL: fname, Checksum: checksum}
		if err := pkg.parseVersion(name); err == nil {
			pkgs = append(pkgs, pkg)
		} else {
			fmt.Printf("ignore invalid version pkg: %s, %v\n", fname, err)
		}
	}
	return pkgs, nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func (self *LocalGetter) Get(pkg *Package, dir string) (string, error) {
	dstFile := path.Join(dir, pkg.Name)
	chkSum, _ := md5sumFile(dstFile)
	if chkSum != "" && chkSum == pkg.Checksum {
		return dstFile, nil
	}
	_, err := copyFile(pkg.URL, dstFile)
	return dstFile, err
}
