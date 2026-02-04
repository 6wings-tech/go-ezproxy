package srv

import (
	"encoding/json"
	"ezp/pkg/git"
	"log"
	"net/http"
)

// $base/$module/@v/list
func modListHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)

	log.Printf("<<🔵 [RepoReq] mod %q VERSIONS LIST: request recv", mod)

	repo, err := git.FindRepo(mod)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q VERSIONS LIST: unable get a repo: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tags := []git.Tag{}
	for _, info := range repo.Tags {
		tags = append(tags, info)
	}

	b, err := json.Marshal(tags)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q VERSIONS LIST: unable to marshal (%d): %v", mod, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)

	log.Printf(">>🟢 [RepoReq] mod %q VERSIONS LIST: resp sent!", mod)
}

// $base/$module/@v/$version.info
func modInfoHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)
	ver := r.Context().Value("ver").(string)

	log.Printf("<<🔵 [RepoReq] mod %q@%s VERSION: request recv", mod, ver)

	_, tag, err := git.FindRepoOfVer(mod, ver)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q@%s VERSION: unable get a repo: %v", mod, ver, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(tag)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q@%s VERSION: unable to marshal (%d): %v", mod, ver, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)

	log.Printf(">>🟢 [RepoReq] mod %q@%s VERSION: ver %s found! resp sent!", mod, ver, ver)
}

// $base/$module/@v/$version.mod
func modModHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)
	ver := r.Context().Value("ver").(string)

	log.Printf("<<🔵 [RepoReq] mod %q@%s GO.MOD file: request recv", mod, ver)

	repo, tag, err := git.FindRepoOfVer(mod, ver)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q@%s GO.MOD file: unable get a repo: %v", mod, err, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, err := repo.GoModFileContent(tag)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q@%s GO.MOD file: unable get go.mod (%d): %v", mod, ver, http.StatusNotFound, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(file))

	log.Printf(">>🟢 [RepoReq] mod %q@%s GO.MOD file: go.mod found! resp sent!", mod, ver)
}

// $base/$module/@v/$version.zip
func modZipHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)
	ver := r.Context().Value("ver").(string)

	log.Printf("<<🔵 [RepoReq] mod %q@%s ZIP: request recv", mod, ver)

	repo, tag, err := git.FindRepoOfVer(mod, ver)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q@%s ZIP: unable get a repo: %v", mod, ver, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// git archive --format=zip --output /tmp/go-swlog.master.zip master
	zipFile, err := repo.CreateZipArchive(tag)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q@%s ZIP: unable create zip archive (%d)", mod, ver, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, zipFile)

	log.Printf(">>🟢 [RepoReq] mod %q@%s ZIP: file served! resp sent!", mod, ver)
}

// $base/$module/@v/?go-get=1
func modGoImportHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)

	log.Printf("<<🔵 [RepoReq] mod %q ?GO-GET=1: request recv", mod)

	repo, err := git.FindRepo(mod)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q ?GO-GET=1: unable get a repo: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s := repo.ImportHtml()
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(s))

	log.Printf(">>🟢 [RepoReq] mod %q ?GO-GET=1: repo found! resp sent!", mod)
}

// $base/$module/@latest
func modLatestHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)

	log.Printf("<<🔵 [RepoReq] mod %q @LATEST: request recv", mod)

	_, tag, err := git.FindRepoWithLatestVer(mod)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q @LATEST: unable get a repo: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(tag)
	if err != nil {
		log.Printf(">>🔴 [RepoReq] mod %q @LATEST: unable to marshal (%d): %v", mod, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)

	log.Printf(">>🟢 [RepoReq] mod %q @LATEST: repo&lastTag found! resp sent!", mod)
}
