package srv

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type server struct {
	// host:port
	addr string
}

func (s *server) Run() {
	mux := s.mux()

	srv := http.Server{
		Addr:         s.addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      mux,
	}

	log.Printf("[SRV] Starting a server at %s", s.addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("[SRV] %v", err)
	}
}

// https://go.dev/ref/mod#goproxy-protocol
func (s *server) mux() *http.ServeMux {
	mux := http.NewServeMux()

	reModDetPath := regexp.MustCompile(`(.+?)\/@v\/(v[0-9]+\.[0-9]+\.[0-9]+)\.(info|mod|zip)$`)
	reOtherPath := regexp.MustCompile(`(.+?)\/@v\/.*?`)
	reGoGetPath := regexp.MustCompile(`(.+?)\/*\?go-get=1`)

	fixMod := func(mod string) string {
		if mod[0] == byte('/') {
			mod = string([]byte(mod)[1:])
		}
		return mod
	}

	// Handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := ""

		matches := reModDetPath.FindAllStringSubmatch(r.URL.Path, -1)
		mod := ""
		ver := ""

		// Mod version request
		if len(matches) > 0 {
			mod = fixMod(matches[0][1])
			ver = matches[0][2]
			path = matches[0][3]

			newCtx := context.WithValue(r.Context(), "mod", mod)
			newCtx = context.WithValue(newCtx, "ver", ver)

			r = r.WithContext(newCtx)

			switch path {
			// $base/$module/@v/$version.info
			case "info":
				modInfoHandler(w, r)

			// $base/$module/@v/$version.mod
			case "mod":
				modModHandler(w, r)

			// $base/$module/@v/$version.zip
			case "zip":
				modZipHandler(w, r)
			}

			return
		}

		// Other mod request (/@v/list, /@v/@latest)
		matches = reOtherPath.FindAllStringSubmatch(r.URL.Path, -1)
		if len(matches) > 0 {
			switch {
			// mod list
			// $base/$module/@v/list
			case strings.HasSuffix(r.URL.Path, "/@v/list"):
				path = "list"

			// mod latest
			// $base/$module/@latest
			case strings.HasSuffix(r.URL.Path, "/@latest"):
				path = "latest"
			}

			mod = fixMod(matches[0][1])
			newCtx := context.WithValue(r.Context(), "mod", mod)

			r = r.WithContext(newCtx)
		}

		switch path {
		// mod list
		// $base/$module/@v/list
		case "list":
			modListHandler(w, r)
			return

		// mod latest
		// $base/$module/@latest
		case "latest":
			modLatestHandler(w, r)
			return
		}

		// $base/$module?go-get=1
		matches = reGoGetPath.FindAllStringSubmatch(r.URL.String(), -1)
		if len(matches) > 0 {
			mod = fixMod(matches[0][1])
			newCtx := context.WithValue(r.Context(), "mod", mod)

			r = r.WithContext(newCtx)

			modGoImportHandler(w, r)
			return
		}

		http.NotFound(w, r)
	})

	return mux
}

func (s *server) Stop() {
}
