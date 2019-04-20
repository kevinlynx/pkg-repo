package pkgrepo

import (
	"fmt"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestMatchLatest(t *testing.T) {
	assert := assert.New(t)
	pkgs := []*Package{
		&Package{Name: "p-0.1.0-1.tgz"},
		&Package{Name: "p-0.1.9-1.tgz"},
		&Package{Name: "p-0.1.10-124.tgz"},
		&Package{Name: "p-0.1.10-123.tgz"},
		// invalid: p-0.1.10_123.tgz
	}

	for _, pkg := range pkgs {
		assert.Nil(pkg.parseVersion("p"))
		fmt.Printf("pkg (%s) version (%s)\n", pkg.Name, pkg.ver.String())
	}
	latest := MatchLatest(pkgs)
	assert.NotNil(latest)
	fmt.Printf("latest: %s\n", latest.Name)
	assert.Equal("p-0.1.10-124.tgz", latest.Name)
}
