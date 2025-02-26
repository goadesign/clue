package generate

import (
	"fmt"
	"math"
	"path"
	"strings"
)

type (
	Import interface {
		PkgName() string
		PkgPath() string
		Alias() string
		AliasOrPkgName() string
	}

	importMap map[string]Import

	importImpl struct {
		pkgPath, pkgName, alias string
	}
)

func newImport(pkgPath string, args ...string) Import {
	var pkgName, alias string
	switch len(args) {
	case 0:
		pkgName = path.Base(pkgPath)
	case 1:
		pkgName = args[0]
		if pkgName != path.Base(pkgPath) {
			alias = pkgName
		}
	case 2:
		pkgName = args[0]
		alias = args[1]
	}
	return &importImpl{pkgPath: pkgPath, pkgName: pkgName, alias: alias}
}

func (i *importImpl) PkgName() string {
	return i.pkgName
}

func (i *importImpl) PkgPath() string {
	return i.pkgPath
}

func (i *importImpl) Alias() string {
	return i.alias
}

func (i *importImpl) AliasOrPkgName() string {
	if i.alias != "" {
		return i.alias
	}
	return i.pkgName
}

func addImport(pkgImport Import, stdImports, extImports, intImports importMap, modPath string) Import {
	var (
		allImports = []map[string]Import{stdImports, extImports, intImports}
		pkgName    = pkgImport.PkgName()
		pkgPath    = pkgImport.PkgPath()
		dup        = false
	)
	for _, imports := range allImports {
		if i, ok := imports[pkgName]; ok {
			if i.PkgPath() == pkgPath {
				return i
			}
			dup = true
		}
	}
	alias := pkgName
	if dup {
	aliases:
		for i := 1; i <= math.MaxInt; i++ {
			alias = fmt.Sprintf("%v%v", pkgName, i)
			for _, imports := range allImports {
				if i, ok := imports[alias]; ok {
					if i.PkgPath() == pkgPath {
						return i
					}
					continue aliases
				}
			}
			pkgImport = newImport(pkgPath, pkgName, alias)
			break
		}
	}
	if !strings.Contains(pkgPath, ".") {
		stdImports[alias] = pkgImport
	} else if pkgPath == modPath || strings.HasPrefix(pkgPath, modPath+"/") {
		intImports[alias] = pkgImport
	} else {
		extImports[alias] = pkgImport
	}
	return pkgImport
}
