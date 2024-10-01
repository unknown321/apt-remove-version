package reader

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"slices"
)

type Package struct {
	Name    string
	Version string
	Data    string
}

type Packages struct {
	FileName string
	modified bool
	packages []Package
}

var (
	PackagePrefix = []byte("Package: ")
	VersionPrefix = []byte("Version: ")
)

func (ps *Packages) Append(p Package) {
	ps.packages = append(ps.packages, p)
}

func (ps *Packages) Get(name string) []Package {
	res := []Package{}
	for _, p := range ps.packages {
		if p.Name == name {
			res = append(res, p)
		}
	}

	return res
}

func (ps *Packages) Remove(name string, version string) {
	ps.packages = slices.DeleteFunc(ps.packages, func(p Package) bool {
		r, err := regexp.Compile(name)
		if err != nil {
			return false
		}

		if r.Match([]byte(p.Name)) && p.Version == version {
			ps.modified = true
			slog.Info("removing package", "name", p.Name, "version", p.Version, "filename", ps.FileName)
			return true
		}

		return false
	})
}

func (ps *Packages) Save(filename string) error {
	if !ps.modified {
		slog.Debug("not modified", "filename", ps.FileName)
		return nil
	}
	slog.Info("saving", "to", filename, "from", ps.FileName)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", filename, err)
	}

	for _, p := range ps.packages {
		f.Write(PackagePrefix)
		f.Write([]byte(p.Name))
		f.Write([]byte("\n"))
		f.Write(VersionPrefix)
		f.Write([]byte(p.Version))
		f.Write([]byte("\n"))
		f.Write([]byte(p.Data))
		f.Write([]byte("\n"))
	}

	return f.Close()
}

func (ps *Packages) Len() int {
	return len(ps.packages)
}

func Open(filename string) (*Packages, error) {
	var err error
	var data []byte
	if data, err = os.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", filename, err)
	}

	lines := bytes.Split(data, []byte("\n"))

	packages := Packages{FileName: filename}

	var ready bool
	p := Package{}
	for _, line := range lines {
		if bytes.Equal(line, nil) {
			if ready && len(p.Name) > 0 && len(p.Version) > 0 {
				packages.Append(p)
				p = Package{}
			}
			ready = false
			continue
		}

		if bytes.HasPrefix(line, PackagePrefix) {
			if len(p.Name) == 0 {
				p.Name = string(bytes.TrimPrefix(line, PackagePrefix))
			} else {
				return nil, fmt.Errorf("parsing package name %s while previous package %s is not parsed", string(bytes.TrimPrefix(line, PackagePrefix)), p.Name)
			}
		}

		if bytes.HasPrefix(line, VersionPrefix) {
			if len(p.Version) == 0 {
				p.Version = string(bytes.TrimPrefix(line, VersionPrefix))
			} else {
				return nil, fmt.Errorf("parsing package version %s while previous package %s=%s is not parsed", string(bytes.TrimPrefix(line, VersionPrefix)), p.Name, p.Version)
			}
		}

		if ready {
			p.Data += string(line) + "\n"
		}

		if len(p.Name) > 0 && len(p.Version) > 0 {
			ready = true
		}
	}

	return &packages, nil
}
