package pkgrepo

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/ini.v1"
	"net/url"
	"os/user"
	"path"
	"strings"
)

type OssConf struct {
	Host      string `ini:"host"`
	AccessId  string `ini:"accessid"`
	AccessKey string `ini:"accesskey"`
}

type OssGetter struct {
	client   *oss.Client
	bucket   *oss.Bucket
	rootPath string
}

func loadConf(file string) (*OssConf, error) {
	if file == "" {
		if curUser, err := user.Current(); err != nil {
			return nil, err
		} else {
			file = path.Join(curUser.HomeDir, ".osscredentials")
		}
	}
	cfg, err := ini.Load(file)
	if err != nil {
		return nil, err
	}
	var conf OssConf
	if err = cfg.Section("OSSCredentials").MapTo(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func NewOssGetter(url0 string, confile string) (*OssGetter, error) {
	u, err := url.Parse(url0)
	if err != nil {
		return nil, err
	}
	conf, err := loadConf(confile)
	if err != nil {
		return nil, err
	}
	client, err := oss.New(conf.Host, conf.AccessId, conf.AccessKey)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(u.Host)
	if err != nil {
		return nil, err
	}
	return &OssGetter{
		client: client, bucket: bucket, rootPath: strings.TrimPrefix(u.Path, "/"),
	}, nil
}

func (self *OssGetter) List(pattern string) ([]*Package, error) {
	lor, err := self.bucket.ListObjects(oss.Prefix(path.Join(self.rootPath, pattern)))
	if err != nil {
		return nil, err
	}
	pkgs := make([]*Package, 0, len(lor.Objects))
	for _, object := range lor.Objects {
		_, name := path.Split(object.Key)
		pkg := &Package{Name: name, URL: object.Key, Checksum: strings.Trim(object.ETag, "\"")}
		if err := pkg.parseVersion(pattern); err == nil {
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs, nil
}

func (self *OssGetter) Get(pkg *Package, dir string) (string, error) {
	dstFile := path.Join(dir, pkg.Name)
	chkSum, _ := md5sumFile(dstFile)
	if strings.ToUpper(chkSum) == pkg.Checksum {
		return dstFile, ErrCached
	}
	err := self.bucket.GetObjectToFile(pkg.URL, dstFile)
	return dstFile, err
}
