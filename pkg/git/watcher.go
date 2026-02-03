package git

import (
	"bytes"
	"errors"
	"ezp/pkg/cfg"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

func ReposWatcher() {
	dur, err := time.ParseDuration(cfg.C.RescanTmo)
	if err != nil {
		log.Fatalf("[ERR ] bad rescan timeout: %v", err)
	}

	tkr := time.NewTicker(time.Millisecond) // Immediately scanning at the first launch
	defer tkr.Stop()

	for range tkr.C {
		rs.mu.Lock()
		tkr.Stop()

		log.Printf("[INFO] Rescaning repos each %s", dur.String())

		for _, gitDir := range cfg.C.GitDirs {
			entries, err := os.ReadDir(gitDir)
			if err != nil {
				log.Fatal(err)
			}

			for _, e := range entries {
				if !e.IsDir() {
					continue
				}

				// Search for "/.git"
				repoDir := gitDir + "/" + e.Name()

				fi, err := os.Stat(repoDir)
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						log.Printf("[WARN] skip repo %q: no git repo was found inside it", repoDir)
						continue
					} else {
						log.Fatalf("[ERR ] %v", err)
					}
				} else if !fi.IsDir() {
					log.Printf("[WARN] skip repo %q: is not a git repo", repoDir)
					continue
				}

				tags, latest, err := getRepoTags(repoDir)
				if err != nil {
					log.Printf("[WARN] skip repo %q: unable get tags: %v", repoDir, err)
					continue
				}

				if len(tags) == 0 {
					log.Printf("[WARN] skip repo %q: no tags were found", repoDir)
					continue
				}

				mod := e.Name()

				r, ok := rs.repos[mod]
				if ok && r.RepoDir != repoDir {
					log.Printf("[WARN] skip repo %q: a repo with the same name (%q) was already found in %q", repoDir, mod, r.RepoDir)
					continue
				}

				rs.repos[mod] = repo{
					Mod:       mod,
					ModUrl:    cfg.C.Domain + "/" + mod,
					CloneUrl:  cfg.C.CloneUrlPrefix + repoDir, // Example, ssh://git@g.6wings.tech:12345/var/git/go-swlog
					RepoDir:   repoDir,
					Tags:      tags,
					LatestTag: latest,
				}

				log.Printf("[ OK ] repo %q was loaded with %d tag(s)", repoDir, len(rs.repos[mod].Tags))
			}
		}

		log.Printf("[INFO] %d repo(s) loaded", len(rs.repos))

		tkr.Reset(dur)
		rs.mu.Unlock()
	}
}

func getRepoTags(repoDir string) (map[verT]Tag, Tag, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return nil, Tag{}, err
	}

	defer func() {
		if err := os.Chdir(curDir); err != nil {
			log.Fatalf("unable change back to dir %q: %v", curDir, err)
		}
	}()

	err = os.Chdir(repoDir)
	if err != nil {
		return nil, Tag{}, err
	}

	cmd := exec.Command("git", "tag")

	var stdOut bytes.Buffer

	multi := io.MultiWriter(os.Stdout, &stdOut)
	cmd.Stdout = multi

	if err := cmd.Run(); err != nil {
		return nil, Tag{}, err
	}

	var list []string

	tokens := strings.Split(stdOut.String(), "\n")
	for _, token := range tokens {
		if semver.IsValid(token) {
			list = append(list, token)
		}
	}

	tags := make(map[verT]Tag)
	var latestTag Tag

	for _, ver := range list {
		tags[ver] = Tag{
			Version: ver,
			Time:    time.Now(),
		}
	}

	latest := ""
	if len(list) > 0 {
		latest = list[len(list)-1]
		latestTag.Version = latest
		latestTag.Time = time.Now()
	}

	return tags, latestTag, nil
}
