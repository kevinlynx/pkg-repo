package pkgrepo

import (
    "testing"
    "fmt"
    assert "github.com/stretchr/testify/require"
)

func TestOssMatch(t *testing.T) {
    assert := assert.New(t)
    getter, err := NewOssGetter("oss://drogopub/tmp", "")
    assert.Nil(err)
    assert.Equal("tmp", getter.rootPath)
    pkgs, err := getter.Match("carbon")
    assert.Nil(err)
    for _, pkg := range pkgs {
        fmt.Printf("%+v\n", pkg)
    }
    if len(pkgs) > 0 {
        _, err := getter.Get(pkgs[0], "/tmp")
        assert.Nil(err)
    }
}
