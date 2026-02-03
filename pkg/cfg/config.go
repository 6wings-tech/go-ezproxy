package cfg

type config struct {
	GitDirs        []string `json:"gitdirs"`
	RescanTmo      string   `json:"rescanTmo"`
	Domain         string   `json:"domain"`
	CloneUrlPrefix string   `json:"cloneUrlPrefix"`
}
