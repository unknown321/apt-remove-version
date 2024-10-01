package cmd

import (
	"apt-remove-version/reader"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Run() error {
	var listsPath string
	var outPath string
	var pkgPath string
	var keepCache bool
	flag.StringVar(&listsPath, "lists", "/var/lib/apt/lists/", "path to apt lists directory with *_Packages")
	flag.StringVar(&outPath, "out", "/var/lib/apt/lists/", "path to output directory")
	flag.StringVar(&pkgPath, "pkgcache", "/var/cache/apt/pkgcache.bin", "full path to apt cache file")
	flag.BoolVar(&keepCache, "keep_cache", false, "keep apt cache file?")
	out := os.Stdout
	flag.CommandLine.SetOutput(out)
	flag.Usage = func() {
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(out, "\nPackage name must be provided in RE2 (https://golang.org/s/re2syntax) format.\n")
		fmt.Fprint(out, "Package version must be an exact match.\n")
		fmt.Fprint(out, "\nExample:\n")
		fmt.Fprint(out, "\t./apt-remove-version -out /tmp/no-new-nvidia/ nvidia-driver=550.90.12-1 \".*=545.23.08-1\"\n")
	}
	flag.Parse()

	rem := flag.Args()
	if len(rem) == 0 {
		return fmt.Errorf("no packages to remove provided")
	}

	toRemove := []reader.Package{}
	for _, v := range rem {
		res := strings.Split(v, "=")
		if len(res) != 2 {
			continue
		}

		toRemove = append(toRemove, reader.Package{Name: res[0], Version: res[1]})
	}

	if len(toRemove) == 0 {
		return fmt.Errorf("no packages to remove provided")
	}

	packageFiles, err := filepath.Glob(listsPath + "*_Packages")
	if err != nil {
		return fmt.Errorf("cannot list apt lists directory: %w", err)
	}

	for _, packageFile := range packageFiles {
		packages, err := reader.Open(packageFile)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", packageFile, err)
		}

		for _, p := range toRemove {
			packages.Remove(p.Name, p.Version)
		}

		out := path.Join(outPath, path.Base(packageFile))
		if err = packages.Save(out); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	if !keepCache {
		if err = os.Remove(pkgPath); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("cannot remove apt cache file: %w", err)
			}
		} else {
			slog.Info("removed apt cache file", "path", pkgPath)
		}
	}

	return nil
}
