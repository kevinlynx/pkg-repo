package pkgrepo

type Package struct {
    Name string
    URL string
    Checksum string
}

type Getter interface {
    Get(pkg *Package, path string) (string, error)
    Match(pattern string) ([]*Package, error)
}

