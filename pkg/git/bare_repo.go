package git

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type bareRepo struct {
	ModifTime time.Time
}

func NewBareRepo(repoRoot string) (bareRepo, error) {
	rfi, err := os.Stat(repoRoot)
	if err != nil {
		return bareRepo{}, err
	} else if !rfi.IsDir() {
		return bareRepo{}, fmt.Errorf("%q not a dir", repoRoot)
	}

	br := bareRepo{
		ModifTime: rfi.ModTime(),
	}

	// Bare repo standard dirs & files
	repoStdRes := []string{
		"HEAD",
		"branches/",
		"config",
		"description",
		"hooks/",
		"info/",
		"objects/",
		"refs/",
	}

	for _, res := range repoStdRes {
		expDir := string([]byte(res)[len(res)-1:]) == "/"

		res = repoRoot + "/" + res
		fi, err := os.Stat(res)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				resType := "file"
				if expDir {
					resType = "dir"
				}
				return bareRepo{}, fmt.Errorf("%s %q not found. Maybe %q it's not a bare repo?", resType, res, repoRoot)
			} else {
				return bareRepo{}, err
			}
		} else if expDir && !fi.IsDir() {
			return bareRepo{}, fmt.Errorf("res %q is not a dir of a bare repo", res)
		} else if !expDir && !fi.Mode().IsRegular() {
			return bareRepo{}, fmt.Errorf("res %q is not a file of a bare repo", res)
		}

		if strings.HasSuffix(res, "objects/") && fi.ModTime().After(br.ModifTime) {
			br.ModifTime = fi.ModTime()
		}
	}

	return br, nil
}
