package pkgrepo

import (
	"fmt"
	assert "github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const (
	testFilePath     = "/tmp/pkgrepo"
	testPackageCache = "/tmp/pkgrepo_cache"
)

func writeTestFiles() {
	os.MkdirAll(testFilePath, 0755)
	for i := 0; i < 3; i++ {
		s := fmt.Sprintf("hello%d", i)
		ioutil.WriteFile(path.Join(testFilePath, fmt.Sprintf("my-pkg-0.1.%d-%d.tgz", i, i)), []byte(s), 0755)
	}
}

func TestLocalGetter(t *testing.T) {
	assert := assert.New(t)
	writeTestFiles()
	getter, _ := NewGetter("file://" + testFilePath)
	pkgs, err := getter.List("my-pkg", "")
	assert.Nil(err)
	for _, pkg := range pkgs {
		fmt.Printf("%+v\n", *pkg)
	}
	latest := MatchLatest(pkgs)
	os.MkdirAll(testPackageCache, 0755)
	_, err = getter.Get(latest, testPackageCache)
	assert.Nil(err)
}
