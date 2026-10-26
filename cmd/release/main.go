// Command release produces self-contained archives and checksums without external packagers.
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	version := flag.String("version", "1.2.0", "release version")
	output := flag.String("output", "dist", "output directory")
	flag.Parse()
	if *version == "" || strings.ContainsAny(*version, "/\\ \t\r\n") {
		panic("invalid version")
	}
	if err := build(*version, *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func build(version, output string) error {
	if err := os.MkdirAll(output, 0755); err != nil {
		return err
	}
	temp, err := os.MkdirTemp("", "menucraft-release-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	commit := "unknown"
	if data, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
		commit = strings.TrimSpace(string(data))
	}
	date := time.Now().UTC().Format(time.RFC3339)
	licenses, err := dependencyLicenses()
	if err != nil {
		return err
	}
	flags := fmt.Sprintf("-s -w -X main.version=%s -X main.commit=%s -X main.buildDate=%s", version, commit, date)
	var archives []string
	for _, target := range [][2]string{{"windows", "amd64"}, {"windows", "arm64"}, {"linux", "amd64"}, {"linux", "arm64"}} {
		name := "menucraft_" + version + "_" + target[0] + "_" + target[1]
		dir := filepath.Join(temp, name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		suffix := ""
		if target[0] == "windows" {
			suffix = ".exe"
		}
		for _, binary := range [][2]string{{"menucraft", "."}, {"seed", "./cmd/seed"}} {
			command := exec.Command("go", "build", "-trimpath", "-buildvcs=false", "-ldflags", flags, "-o", filepath.Join(dir, binary[0]+suffix), binary[1])
			for _, v := range os.Environ() {
				if !strings.HasPrefix(v, "GOOS=") && !strings.HasPrefix(v, "GOARCH=") && !strings.HasPrefix(v, "CGO_ENABLED=") {
					command.Env = append(command.Env, v)
				}
			}
			command.Env = append(command.Env, "GOOS="+target[0], "GOARCH="+target[1], "CGO_ENABLED=0")
			command.Stdout, command.Stderr = os.Stdout, os.Stderr
			if err := command.Run(); err != nil {
				return err
			}
		}
		for _, file := range []string{"README.md", "README.fa.md", "README.en.md", "CHANGELOG.md", "docs/DEPLOYMENT.md", "docs/WORDPRESS.md", "docs/NOTEBOOK_AND_QR.md", "docs/THIRD_PARTY_NOTICES.md", "docs/ARCHITECTURE.md", "docs/VERIFICATION.md", "docs/screenshots/menu-fa.jpg"} {
			if err := copyFile(file, filepath.Join(dir, file)); err != nil {
				return err
			}
		}
		for name, source := range licenses {
			if err := copyFile(source, filepath.Join(dir, "licenses", name)); err != nil {
				return err
			}
		}
		if target[0] == "windows" {
			script := "@echo off\r\ncd /d \"%~dp0\"\r\nif not defined MENU_ADDR set MENU_ADDR=127.0.0.1:8080\r\necho Open http://127.0.0.1:8080 in your browser.\r\nmenucraft.exe\r\npause\r\n"
			if err := os.WriteFile(filepath.Join(dir, "run.bat"), []byte(script), 0644); err != nil {
				return err
			}
		} else {
			script := "#!/bin/sh\nset -eu\ncd -- \"$(dirname -- \"$0\")\"\nexport MENU_ADDR=\"\u0024{MENU_ADDR:-127.0.0.1:8080}\"\nexec ./menucraft\n"
			if err := os.WriteFile(filepath.Join(dir, "run.sh"), []byte(script), 0755); err != nil {
				return err
			}
		}
		archive := filepath.Join(output, name+".zip")
		if target[0] == "linux" {
			archive = filepath.Join(output, name+".tar.gz")
			err = tarDir(dir, archive)
		} else {
			err = zipDir(dir, archive, name)
		}
		if err != nil {
			return err
		}
		archives = append(archives, archive)
		fmt.Println("Created", archive)
	}
	plugin := filepath.Join(output, "menucraft-wordpress_"+version+".zip")
	if err := zipDir("wordpress/menucraft", plugin, "menucraft"); err != nil {
		return err
	}
	archives = append(archives, plugin)
	sort.Strings(archives)
	var checks strings.Builder
	for _, path := range archives {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		h := sha256.New()
		_, err = io.Copy(h, f)
		f.Close()
		if err != nil {
			return err
		}
		fmt.Fprintf(&checks, "%s  %s\n", hex.EncodeToString(h.Sum(nil)), filepath.Base(path))
	}
	return os.WriteFile(filepath.Join(output, "SHA256SUMS.txt"), []byte(checks.String()), 0644)
}

func dependencyLicenses() (map[string]string, error) {
	data, err := exec.Command("go", "list", "-m", "-json", "all").Output()
	if err != nil {
		return nil, err
	}
	licenses := map[string]string{"Vazirmatn-OFL.txt": "web/static/fonts/OFL.txt"}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	for {
		var module struct {
			Path, Version, Dir string
			Main               bool
		}
		if err := decoder.Decode(&module); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if module.Main || module.Dir == "" {
			continue
		}
		entries, err := os.ReadDir(module.Dir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			name := strings.ToUpper(entry.Name())
			if !entry.IsDir() && (strings.HasPrefix(name, "LICENSE") || strings.HasPrefix(name, "COPYING") || strings.HasPrefix(name, "NOTICE")) {
				licenses[filepath.Join(filepath.FromSlash(module.Path+"@"+module.Version), entry.Name())] = filepath.Join(module.Dir, entry.Name())
			}
		}
	}
	root, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return nil, err
	}
	licenses["Go-LICENSE.txt"] = filepath.Join(strings.TrimSpace(string(root)), "LICENSE")
	return licenses, nil
}
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
func zipDir(dir, out, prefix string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		h.Name = prefix + "/" + filepath.ToSlash(rel)
		h.Method = zip.Deflate
		entry, err := w.CreateHeader(h)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = entry.Write(data)
		return err
	})
	if err != nil {
		_ = w.Close()
		return err
	}
	return w.Close()
}
func tarDir(dir, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		h, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		h.Name = filepath.Base(dir) + "/" + filepath.ToSlash(rel)
		if filepath.Base(path) == "menucraft" || filepath.Base(path) == "seed" || strings.HasSuffix(path, ".sh") {
			h.Mode = 0755
		}
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = tw.Write(data)
		return err
	})
	closeTar := tw.Close()
	closeGzip := gz.Close()
	if err != nil {
		return err
	}
	if closeTar != nil {
		return closeTar
	}
	return closeGzip
}
