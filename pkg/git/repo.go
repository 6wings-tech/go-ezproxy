package git

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Version as semver (x.y.z)
// For example, "1.23.4567"
type verT = string

type repo struct {
	Mod       string
	ModUrl    string
	CloneUrl  string
	RepoDir   string
	Tags      map[verT]Tag // sorted versions
	LatestTag Tag          // latest version
}

// ImportHtml generates the instruction strcing for the "go get" command
func (r repo) ImportHtml() string {
	s := fmt.Sprintf("<meta name=\"go-import\" content=\"%s %s %s\">", r.ModUrl, "git", r.CloneUrl)
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

	err = os.Chdir(r.RepoDir)
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

	err = os.Chdir(r.RepoDir)
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
