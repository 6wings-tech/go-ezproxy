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

	log.Printf("mod info: mod %q LIST requested", mod)

	repo, err := git.FindRepo(mod)
	if err != nil {
		log.Printf("[LIST] unable get %s: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tags := []git.Tag{}
	for _, info := range repo.Tags {
		tags = append(tags, info)
	}

	b, err := json.Marshal(tags)
	if err != nil {
		log.Printf("[LIST] unable marshal %s: %v", mod, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[LIST] reply: %v", string(b))

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// $base/$module/@v/$version.info
func modInfoHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)
	ver := r.Context().Value("ver").(string)

	log.Printf("mod info: mod %q version %s INFO requested", mod, ver)

	_, tag, err := git.FindRepoOfVer(mod, ver)
	if err != nil {
		log.Printf("[VERSION] unable get %s: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(tag)
	if err != nil {
		log.Printf("[VERSION] unable marshal %s: %v", mod, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[VERSION] reply: %v", string(b))

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// $base/$module/@v/$version.mod
func modModHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)
	ver := r.Context().Value("ver").(string)

	log.Printf("mod info: mod %q version %s MOD requested", mod, ver)

	repo, tag, err := git.FindRepoOfVer(mod, ver)
	if err != nil {
		log.Printf("[MOD] unable get %s: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, err := repo.GoModFileContent(tag)
	if err != nil {
		log.Printf("[MOD] unable get go.mod of %s: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("[MOD] reply: %v", file)

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(file))
}

// $base/$module/@v/$version.zip
func modZipHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)
	ver := r.Context().Value("ver").(string)

	log.Printf("mod info: mod %s@%s zip archive requested", mod, ver)

	repo, tag, err := git.FindRepoOfVer(mod, ver)
	if err != nil {
		log.Printf("[ZIP] unable get %s@%s: %v", mod, ver, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// git archive --format=zip --output /tmp/go-swlog.master.zip master
	zipFile, err := repo.CreateZipArchive(tag)
	if err != nil {
		log.Printf("[ZIP] unable create zip archive %s@%s: %v", mod, ver, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, zipFile)
}

// $base/$module/@v/?go-get=1
func modGoImportHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)

	log.Printf("[GO-GET=1] mod info: mod %q Go-GET=1 requested", mod)

	repo, err := git.FindRepo(mod)
	if err != nil {
		log.Printf("[MOD] unable get %s: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("[GO-GET=1] unable get %s: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s := repo.ImportHtml()

	log.Printf("[GO-GET=1] reply: %v", s)

	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(s))
}

// $base/$module/@latest
func modLatestHandler(w http.ResponseWriter, r *http.Request) {
	mod := r.Context().Value("mod").(string)
	log.Printf("mod info: mod %q LATEST requested", mod)

	_, tag, err := git.FindRepoWithLatestVer(mod)
	if err != nil {
		log.Printf("[VERSION] unable get %s: %v", mod, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(tag)
	if err != nil {
		log.Printf("[VERSION] unable marshal %s: %v", mod, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[VERSION] reply: %v", string(b))

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
