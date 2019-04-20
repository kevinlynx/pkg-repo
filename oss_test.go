package pkgrepo

import (
	"fmt"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestOssList(t *testing.T) {
	assert := assert.New(t)
	getter, err := NewOssGetter("oss://drogopub/tmp", "")
	assert.Nil(err)
	pkgs, err := getter.List("hemisconf")
	assert.Nil(err)
	for _, pkg := range pkgs {
		fmt.Printf("%+v\n", pkg)
	}
	if len(pkgs) > 0 {
		_, err := getter.Get(pkgs[0], "/tmp")
		assert.True(err == nil || err == ErrCached)
	}
}
