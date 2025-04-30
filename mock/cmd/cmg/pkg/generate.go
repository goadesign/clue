package cluemockgen

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"goa.design/clue/log"
	"goa.design/clue/mock/cmd/cmg/pkg/generate"
	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

func Generate(ctx context.Context, patterns []string, dir string) error {
	ps, err := parse.LoadPackages(patterns, dir)
	if err != nil {
		log.Error(ctx, err)
		return err
	}

	var errs []error

	for _, p := range ps {
		err = generatePackage(ctx, p)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func generatePackage(ctx context.Context, p parse.Package) error {
	ctx = log.With(ctx, log.KV{K: "pkg name", V: p.Name()})
	log.Print(ctx, log.KV{K: "pkg path", V: p.PkgPath()})

	is, err := p.Interfaces()
	if err != nil {
		log.Error(ctx, err)
		return err
	}

	interfacesByFile := make(map[string][]parse.Interface)
	for _, i := range is {
		ctx := log.With(ctx, log.KV{K: "interface name", V: i.Name()})
		log.Print(ctx, log.KV{K: "is exported", V: i.IsExported()}, log.KV{K: "file", V: i.File()})
		if i.IsExported() {
			var (
				exportedMethods   = 0
				unexportedMethods []string
			)
			for _, method := range i.Methods() {
				if method.IsExported() {
					exportedMethods++
				} else {
					unexportedMethods = append(unexportedMethods, method.Name())
				}
			}
			if len(unexportedMethods) > 0 {
				err = fmt.Errorf("unexported methods: %v", strings.Join(unexportedMethods, ", "))
				log.Error(ctx, err)
				return err
			}
			if exportedMethods <= 0 {
				continue
			}
			interfacesByFile[i.File()] = append(interfacesByFile[i.File()], i)
		}
	}
	for file, interfaces := range interfacesByFile {
		err = generateFile(ctx, p, file, interfaces)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateFile(ctx context.Context, p parse.Package, file string, interfaces []parse.Interface) error {
	ctx = log.With(ctx, log.KV{K: "file", V: file})
	interfaceNames := make([]string, len(interfaces))
	for j, i := range interfaces {
		interfaceNames[j] = i.Name()
	}
	log.Print(ctx, log.KV{K: "interface names", V: interfaceNames})

	dir, baseFile := filepath.Split(file)
	mocksDir := filepath.Join(dir, "mocks")
	mocksFile := filepath.Join(mocksDir, baseFile)
	ctx = log.With(ctx, log.KV{K: "mocks file", V: mocksFile})

	if err := os.MkdirAll(mocksDir, 0o777); err != nil {
		log.Error(ctx, err)
		return err
	}

	f, err := os.CreateTemp(mocksDir, ".*."+baseFile)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	defer func() {
		if removeErr := os.Remove(f.Name()); removeErr != nil && !errors.Is(removeErr, fs.ErrNotExist) {
			log.Error(ctx, fmt.Errorf("failed to remove temporary file: %w", removeErr))
		}
	}()
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Error(ctx, fmt.Errorf("failed to close file: %w", closeErr))
		}
	}()

	mocks := generate.NewMocks("mock", p, interfaces, Version)
	if err := mocks.Render(f); err != nil {
		log.Error(ctx, err)
		return err
	}
	if err := os.Rename(f.Name(), mocksFile); err != nil {
		log.Error(ctx, err)
		return err
	}

	return nil
}
