package pkgrepo

import (
	"fmt"
	"github.com/hashicorp/go-version"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestMatchLatest(t *testing.T) {
	assert := assert.New(t)
	v1, _ := version.NewVersion("0.1.0-201904191552")
	v2, _ := version.NewVersion("0.1.10-1")
	assert.True(v2.GreaterThan(v1))

	pkgs := []*Package{
		&Package{Name: "abc/my-pkg-0.1.9-1.tgz"},
		&Package{Name: "abc/my-pkg-0.1.10-124.tgz"},
		&Package{Name: "abc/my-pkg-0.1.10-1.tgz"},
		&Package{Name: "abc/my-pkg-0.1.0-201904191552.tgz"},
		// invalid: p-0.1.10_123.tgz
	}

	for _, pkg := range pkgs {
		assert.Nil(pkg.parseVersion("my-pkg"))
		fmt.Printf("pkg (%s) version (%s)\n", pkg.Name, pkg.Version())
	}
	latest := MatchLatest(pkgs)
	assert.NotNil(latest)
	fmt.Printf("latest: %s\n", latest.Name)
	assert.Equal("abc/my-pkg-0.1.10-124.tgz", latest.Name)
}
