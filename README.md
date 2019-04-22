Package repo is like a simplified rpm repo, which used to detect and download latest package from remote repo.

```
getter, _ := pkgrepo.NewGetter(repoRoot)
pkgs, _ := getter.List("my-package", "0.1.0")
latestPkg := pkgrepo.MatchLatest(pkgs)
getter.Get(latestPkg, "/local/cache")
```

The remote repo can be a local file system, or an OSS storage on aliyun, or you can write your own remote repo.
