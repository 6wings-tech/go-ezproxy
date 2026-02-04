package git

import (
	"bytes"
	"errors"
	"ezp/pkg/cfg"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

func ReposWatcher() {
	dur, err := time.ParseDuration("15s")
	if err != nil {
		log.Fatalf("🔴 [ERR ] bad rescan timeout: %v", err)
	}

	tkr := time.NewTicker(time.Millisecond) // Immediately scanning at the first launch
	defer tkr.Stop()
	ttlRescans := 0  // a quantity of the repos rescans
	reloaded := 0 // a reload repos quantity (if changes were found) during of the one rescaning

	log.Printf("🔵 [INFO] repos root is %q", cfg.C.ReposRoot)
	log.Printf("🔵 [INFO] checking repos changes planned at each %s", dur.String())

	for range tkr.C {
		rs.mu.Lock()
		tkr.Stop()

		reloaded = 0
		ttlRescans++

		entries, err := os.ReadDir(cfg.C.ReposRoot)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			// Getting info of a BARE repo
			repoRoot := cfg.C.ReposRoot + "/" + e.Name()

			br, err := getBareRepo(repoRoot)
			if err != nil {
				log.Fatalf("🔴 [ERR ] %v", err)
				continue
			}

			mod := e.Name()

			r, ok := rs.repos[mod]
			if ok && r.Root != repoRoot {
				log.Printf("🟡 [WARN] skip repo %q: a repo with the same name (%q) was already found in %q", repoRoot, mod, r.Root)
				continue
			} else if ok {
				bMt := br.ModifTime.In(time.UTC).Format("2006-01-02 15:04:05")
				rMt := r.ModifTime.In(time.UTC).Format("2006-01-02 15:04:05")

				// No changes since of the last commit
				if bMt == rMt {
					continue
				}
			}

			// Tags
			tags, latest, err := getRepoTags(repoRoot)
			if err != nil {
				log.Printf("🟡 [WARN] skip repo %q: unable get tags: %v", repoRoot, err)
				continue
			}

			// Collecting
			rs.repos[mod] = repo{
				Mod:       mod,
				Root:      repoRoot,
				ModifTime: br.ModifTime,
				Tags:      tags,
				LatestTag: latest,
			}

			if len(tags) == 0 {
				log.Printf("🟡 [WARN] repo %q has no tags yet", repoRoot)
			} else {
				log.Printf("🟢 [ OK ] repo %q was loaded with %d tag(s)", repoRoot, len(rs.repos[mod].Tags))
				reloaded++
			}
		}

		if reloaded > 0 {
			log.Printf("🔵 [INFO] %d repo(s) reloaded", len(rs.repos))
		} else if ttlRescans > 0 && ttlRescans%100 == 0 {
			// to prevent disk space exhausting log it each 100 time
			log.Printf("🔵 [INFO] at totally repos in %q rescanned %d time(s)", cfg.C.ReposRoot, ttlRescans)
		}

		tkr.Reset(dur)
		rs.mu.Unlock()
	}
}

func getBareRepo(repoRoot string) (bareRepo, error) {
	fi, err := os.Stat(repoRoot)
	if err != nil {
		return bareRepo{}, err
	} else if !fi.IsDir() {
		return bareRepo{}, fmt.Errorf("%q not a dir", repoRoot)
	}

	// dir "objects/"
	res := repoRoot + "/objects"
	ft, err := os.Stat(res)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return bareRepo{}, fmt.Errorf("subdir %q not found. Maybe %q it's not a bare repo?", res, repoRoot)
		} else {
			return bareRepo{}, err
		}
	} else if !fi.IsDir() {
		return bareRepo{}, fmt.Errorf("subdir %q not found in the bare repo", res)
	}

	br := bareRepo{
		ModifTime: ft.ModTime(),
	}

	return br, nil
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
