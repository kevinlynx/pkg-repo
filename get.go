package pkgrepo

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"io"
	"net/url"
	"os"
	"regexp"
	"sort"
)

var ErrCached = errors.New("package cached")

type Package struct {
	Name     string
	URL      string
	Checksum string
	ver      *version.Version
}

type Getter interface {
	Get(pkg *Package, path string) (string, error)
	List(pattern string) ([]*Package, error)
}

func NewGetter(url0 string, args ...interface{}) (Getter, error) {
	u, err := url.Parse(url0)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "oss" {
		conffile := ""
		if len(args) > 0 {
			if s, ok := args[0].(string); ok {
				conffile = s
			}
		}
		return NewOssGetter(url0, conffile)
	}
	return nil, errors.New("not supported")
}

func MatchLatest(pkgs []*Package) *Package {
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ver.GreaterThan(pkgs[j].ver)
	})
	return pkgs[0]
}

func (pkg *Package) parseVersion(prefix string) error {
	re, err := regexp.Compile(prefix + `[-_](.+)\..+`)
	if err != nil {
		return err
	}
	m := re.FindStringSubmatch(pkg.Name)
	if m == nil || len(m) < 2 {
		return fmt.Errorf("not found version string in package: %s", pkg.Name)
	}
	pkg.ver, err = version.NewVersion(m[1])
	if err != nil {
		return err
	}
	return nil
}

func md5sumFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	m := md5.New()
	if _, err := io.Copy(m, file); err != nil {
		return "", err
	}
	bytes := m.Sum(nil)
	md5sum := fmt.Sprintf("%x", bytes)
	return md5sum, nil
}
