package git

import (
	"bytes"
	"ezp/pkg/cfg"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Version as semver (x.y.z)
// For example, "1.23.4567"
type verT = string

type repo struct {
	Mod       string
	ModifTime time.Time
	Root      string
	Tags      map[verT]Tag
	LatestTag Tag // latest version
}

// CloneUrl returns an ssh url for the `go get` command
// An exmaple of a returned value is "ssh://git@example.com:12345/var/git"
func (r repo) CloneUrl() string {
	host := cfg.C.Domain
	if cfg.C.SshPort != 22 {
		host += ":" + strconv.Itoa(cfg.C.SshPort)
	}
	return "ssh://" + cfg.C.SshUser + "@" + host + r.Root
}

// ModUrl returns a module url that is required by `go get` command
// An example of a returned value for repo `go-myrepo` is "g.6wings.tech/go-myrepo"
// On `go get` or `go tidy` Go will be request "httpS://g.6wings.tech/go-myrepo"
// to check a mod availability and it props
func (r repo) ModUrl() string {
	return cfg.C.Domain + "/" + r.Mod
}

// GoImportHtml generates the instruction strcing for the "go get" command
func (r repo) GoImportHtml() string {
	s := fmt.Sprintf("<meta name=\"go-import\" content=\"%s %s %s\">", r.ModUrl(), "git", r.CloneUrl())
	return s
}

// go.mod content
func (r repo) GoModFileContent(t Tag) (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	defer func() {
		if err := os.Chdir(curDir); err != nil {
			log.Fatalf("unable change back to dir %q: %v", curDir, err)
		}
	}()

	err = os.Chdir(r.Root)
	if err != nil {
		return "", err
	}

	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer

	// git show v1.23.4567:go.mod
	cmd := exec.Command("git", "show", t.Version+":go.mod")
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf

	if err := cmd.Run(); err != nil {
		return "", err
	} else if len(stdErrBuf.Bytes()) > 0 {
		return "", fmt.Errorf("error on getting go.mod of module %s@%s: %s", r.Mod, t.Version, stdErrBuf.String())
	}

	return stdOutBuf.String(), nil
}

func (r repo) CreateZipArchive(t Tag) (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	defer func() {
		if err := os.Chdir(curDir); err != nil {
			log.Fatalf("unable change back to dir %q: %v", curDir, err)
		}
	}()

	err = os.Chdir(r.Root)
	if err != nil {
		return "", err
	}

	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer

	// git archive --format=zip --output=/tmp/go-swlog.v1.0.1.zip v1.0.1             - DONE
	// git archive --format=zip --output=/tmp/go-swlog.master.zip master             - TODO
	file := "/tmp/" + r.Mod + "." + t.Version + ".zip"

	// If an archive is already created
	fi, err := os.Stat(file)
	if err == nil && fi.Mode().IsRegular() {
		return file, err
	}

	cmd := exec.Command("git", "archive", "--format=zip", "--output="+file, t.Version)
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf

	if err := cmd.Run(); err != nil {
		return "", err
	} else if len(stdErrBuf.Bytes()) > 0 {
		return "", fmt.Errorf("error on creating zip archive of module %s@%s: %s", r.Mod, t.Version, stdErrBuf.String())
	}

	if fi, err := os.Stat(file); err != nil {
		return "", err
	} else if !fi.Mode().IsRegular() {
		return "", fmt.Errorf("archive file %s is not a regular file", file)
	}

	return file, nil
}
