package git

import (
	"errors"
	"sync"
)

var rs = &reposStorage{
	repos: make(map[modT]repo),
}

// Mod (repo) name
// For example, "go-my-pet-project" that will be a part of module url:
// https://yourdomain/go-my-pet-project
type modT = string

type reposStorage struct {
	mu    sync.RWMutex
	repos map[modT]repo
}

func FindRepo(mod modT) (repo, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	r, ok := rs.repos[mod]
	if !ok {
		return repo{}, errors.New("repo not found")
	}

	return r, nil
}

func FindRepoOfVer(mod modT, ver verT) (repo, Tag, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	r, ok := rs.repos[mod]
	if !ok {
		return repo{}, Tag{}, errors.New("repo not found")
	}

	if len(r.Tags) == 0 {
		return repo{}, Tag{}, errors.New("repo has no versions commited")
	}

	t, ok := r.Tags[ver]
	if !ok {
		return repo{}, Tag{}, errors.New("version not found")
	}

	return r, t, nil
}

func FindRepoWithLatestVer(mod modT) (repo, Tag, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	r, ok := rs.repos[mod]
	if !ok {
		return repo{}, Tag{}, errors.New("repo not found")
	}

	return r, r.LatestTag, nil
}
